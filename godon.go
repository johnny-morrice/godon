package godon

import (
	"context"
	"net/url"
	"strings"

	"github.com/mattn/go-mastodon"
	"github.com/pkg/errors"
)

type Options struct {
	ServerUrl    *url.URL
	App          *mastodon.Application
	AppConfig    *mastodon.AppConfig
	RefreshToken string
	Context      context.Context
}

type Godon struct {
	Options
}

func New(options Options) (*Godon, error) {
	godon := &Godon{
		Options: options,
	}

	if godon.Context == nil {
		godon.Context = context.Background()
	}

	if godon.ServerUrl == nil && godon.AppConfig.Server == "" {
		return nil, errors.New("ServerUrl missing")
	}

	if godon.AppConfig.Server == "" {
		godon.AppConfig.Server = godon.ServerUrl.String()
	}

	if godon.ServerUrl == nil && godon.AppConfig.Server != "" {
		url, err := url.Parse(godon.AppConfig.Server)

		if err != nil {
			return nil, err
		}

		godon.ServerUrl = url
	}

	return godon, nil
}

func (godon *Godon) Register() error {
	const errMsg = "Godon.Register failed"

	app, err := mastodon.RegisterApp(godon.Context, godon.AppConfig)

	if err != nil {
		return errors.Wrap(err, errMsg)
	}

	godon.App = app
	return nil
}

func (godon *Godon) Authorize(method func(url string) (string, error)) error {
	const errMsg = "Authorize failed"

	const responseTypeParam = "response_type=code"
	const redirectUriParam = "redirect_uri=urn:ietf:wg:oauth:2.0:oob"
	scopesParam := "scopes=" + godon.AppConfig.Scopes
	clientIdParam := "client_id=" + godon.App.ClientID

	url := godon.makeGetUrl("/oauth/authorize", scopesParam, responseTypeParam, redirectUriParam, clientIdParam)
	refreshToken, err := method(url)

	if err != nil {
		return errors.Wrap(err, refreshToken)
	}

	godon.RefreshToken = refreshToken
	return nil
}

func (godon *Godon) makeGetUrl(path string, params ...string) string {
	var getUrl *url.URL
	*getUrl = *godon.ServerUrl
	getUrl.Path = path

	encodedParams := make([]string, len(params))

	for i, p := range params {
		encodedParams[i] = url.QueryEscape(p)
	}
	getUrl.RawQuery = strings.Join(encodedParams, "&")

	return getUrl.String()
}

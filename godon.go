package godon

import (
	"context"
	"net/url"
	"strings"

	"github.com/mattn/go-mastodon"
	"github.com/pkg/errors"
)

type Options struct {
	App          *mastodon.Application
	AppConfig    *mastodon.AppConfig
	RefreshToken string
	Context      context.Context
}

type Godon struct {
	Options
	serverUrl *url.URL
}

func New(options Options) (*Godon, error) {
	godon := &Godon{
		Options: options,
	}

	if godon.Context == nil {
		godon.Context = context.Background()
	}

	if godon.AppConfig.Server == "" {
		return nil, errors.New("Server URL missing")
	}

	url, err := url.Parse(godon.AppConfig.Server)

	if err != nil {
		return nil, err
	}

	godon.serverUrl = url

	godon.AppConfig.RedirectURIs = __REDIRECT_URL_VALUE

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

	query := map[string]string{}
	query[__RESPONSE_TYPE_KEY] = __RESPONSE_TYPE_VALUE
	query[__REDIRECT_URL_KEY] = __REDIRECT_URL_VALUE
	query[__SCOPES_KEY] = godon.AppConfig.Scopes
	query[__CLIENT_ID_KEY] = godon.App.ClientID

	url := godon.makeGetUrl("/oauth/authorize", query)
	refreshToken, err := method(url)

	if err != nil {
		return errors.Wrap(err, refreshToken)
	}

	godon.RefreshToken = refreshToken
	return nil
}

func (godon *Godon) makeGetUrl(path string, query map[string]string) string {
	getUrl := &url.URL{}
	*getUrl = *godon.serverUrl
	getUrl.Path = path

	queryList := make([]string, len(query))

	i := 0
	for k, q := range query {
		queryList[i] = k + "=" + url.QueryEscape(q)
		i++
	}

	getUrl.RawQuery = strings.Join(queryList, "&")

	return getUrl.String()
}

func (godon *Godon) OauthLogin() error {
	panic("not implemented")
}

func (godon *Godon) PostStatus(text string) error {
	panic("not implemented")
}

const __REDIRECT_URL_VALUE = "urn:ietf:wg:oauth:2.0:oob"
const __REDIRECT_URL_KEY = "redirect_uri"
const __SCOPES_KEY = "scopes"
const __CLIENT_ID_KEY = "client_id"
const __RESPONSE_TYPE_KEY = "response_type"
const __RESPONSE_TYPE_VALUE = "code"

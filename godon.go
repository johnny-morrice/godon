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
	AccessToken  string
	Context      context.Context
}

type Godon struct {
	Options
	serverUrl *url.URL
	client    *mastodon.Client
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

type AuthorizationMethod func(url string) (string, error)

func (godon *Godon) Authorize(method AuthorizationMethod) error {
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

func (godon *Godon) OauthLogin(method AuthorizationMethod) error {
	const errMsg = "godon.OauthLogin failed"

	err := godon.Register()

	if err != nil {
		return errors.Wrap(err, errMsg)
	}

	err = godon.Authorize(method)

	if err != nil {
		return errors.Wrap(err, errMsg)
	}

	err = godon.AccessMastodon()

	if err != nil {
		return errors.Wrap(err, errMsg)
	}

	config := godon.makeMastodonConfig()
	godon.client = mastodon.NewClient(config)

	return nil
}

// AccessMastodon fetches an access token and stores within godon.
func (godon *Godon) AccessMastodon() error {
	mastodon.
}

func (godon *Godon) makeMastodonConfig() *mastodon.Config {
	return &mastodon.Config{
		Server:       godon.AppConfig.Server,
		ClientID:     godon.App.ClientID,
		ClientSecret: godon.App.ClientSecret,
		AccessToken:  godon.AccessToken,
	}
}

// TODO should take toot as parameter.
func (godon *Godon) PostStatus(statusText string) error {
	toot := &mastodon.Toot{
		Status: statusText,
	}

	// TODO figure out who should get the returned status.
	_, err := godon.client.PostStatus(godon.Context, toot)

	return err
}

const __REDIRECT_URL_VALUE = "urn:ietf:wg:oauth:2.0:oob"
const __REDIRECT_URL_KEY = "redirect_uri"
const __SCOPES_KEY = "scopes"
const __CLIENT_ID_KEY = "client_id"
const __RESPONSE_TYPE_KEY = "response_type"
const __RESPONSE_TYPE_VALUE = "code"

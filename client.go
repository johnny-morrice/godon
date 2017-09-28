package godon

import (
	"github.com/McKael/madon"
	"github.com/pkg/errors"
)

type ClientOptions struct {
	InstanceName string
	AppName      string
	AppId        string
	AppSecret    string
	WebSite      string
	Token        *madon.UserToken
}

type Client struct {
	ClientOptions
	madonClient *madon.Client
}

func MakeClient(options ClientOptions) *Client {
	return &Client{
		ClientOptions: options,
	}
}

type AuthorizationMethod func(url string) (string, error)

func (client *Client) Login(method AuthorizationMethod) error {
	const errMsg = "godon.Login failed"

	var err error
	if client.Token.AccessToken == "" {
		err = client.register(method)
	} else {
		err = client.restore()
	}

	if err != nil {
		return errors.Wrap(err, errMsg)
	}

	return nil
}

func (client *Client) GetTimeline(timeline string, local bool, limit int) ([]madon.Status, error) {
	const errMsg = "godon.GetTimeline failed"

	lopt := &madon.LimitParams{}

	if limit < 0 {
		lopt.All = true
	} else {
		lopt.Limit = limit
	}

	updates, err := client.madonClient.GetTimelines(timeline, local, lopt)

	if err != nil {
		return nil, errors.Wrap(err, errMsg)
	}

	return updates, nil
}

func (client *Client) authorize(method AuthorizationMethod) error {
	url, err := client.madonClient.LoginOAuth2("", __SCOPES)

	if err != nil {
		return err
	}

	code, err := method(url)

	if err != nil {
		return err
	}

	_, err = client.madonClient.LoginOAuth2(code, __SCOPES)

	return err
}

func (client *Client) register(method AuthorizationMethod) error {
	const errMsg = "godon.register failed"

	mastodonClient, err := madon.NewApp(client.AppName, client.WebSite, __SCOPES, __REDIRECT_URL, client.InstanceName)

	if err != nil {
		return errors.Wrap(err, errMsg)
	}

	client.madonClient = mastodonClient

	err = client.authorize(method)

	if err != nil {
		return errors.Wrap(err, errMsg)
	}

	return nil
}

func (client *Client) restore() error {
	const errMsg = "godon.restore failed"

	mastodonClient, err := madon.RestoreApp(client.AppName, client.InstanceName, client.AppId, client.AppSecret, client.Token)

	if err != nil {
		return errors.Wrap(err, errMsg)
	}

	client.madonClient = mastodonClient
	return nil
}

type Toot struct {
	Status      string
	InReplyTo   int64
	MediaIDs    []int64
	Sensitive   bool
	SpoilerText string
	Visibility  string
}

func (client *Client) GetAppId() string {
	return client.madonClient.ID
}

func (client *Client) GetAppSecret() string {
	return client.madonClient.Secret
}

func (client *Client) PostStatus(toot Toot) error {
	const errMsg = "godon.PostStatus failed"
	// TODO Figure out who to send the status to.
	_, err := client.madonClient.PostStatus(toot.Status, toot.InReplyTo, toot.MediaIDs, toot.Sensitive, toot.SpoilerText, toot.Visibility)

	if err != nil {
		return errors.Wrap(err, errMsg)
	}

	return nil
}

func (client *Client) GetAccessToken() *madon.UserToken {
	if client.madonClient.UserToken == nil {
		return &madon.UserToken{}
	}

	return client.madonClient.UserToken
}

const __REDIRECT_URL = ""

var __SCOPES = []string{"read", "write", "follow"}

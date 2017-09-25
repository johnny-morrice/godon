package godon

import (
	"github.com/McKael/madon"
	"github.com/pkg/errors"
)

type Options struct {
	InstanceName string
	AppName      string
	AppId        string
	AppSecret    string
	WebSite      string
	Token        *madon.UserToken
}

type Godon struct {
	Options
	client *madon.Client
}

func New(options Options) *Godon {
	return &Godon{
		Options: options,
	}
}

type AuthorizationMethod func(url string) (string, error)

func (godon *Godon) Login(method AuthorizationMethod) error {
	const errMsg = "godon.Login failed"

	var err error
	if godon.Token.AccessToken == "" {
		err = godon.register()
	} else {
		err = godon.restore()
	}

	if err != nil {
		return errors.Wrap(err, errMsg)
	}

	err = godon.authorize(method)

	if err != nil {
		return errors.Wrap(err, errMsg)
	}

	return nil
}

func (godon *Godon) authorize(method AuthorizationMethod) error {
	url, err := godon.client.LoginOAuth2("", __SCOPES)

	if err != nil {
		return err
	}

	code, err := method(url)

	if err != nil {
		return err
	}

	_, err = godon.client.LoginOAuth2(code, __SCOPES)

	return err
}

func (godon *Godon) register() error {
	const errMsg = "godon.register failed"

	client, err := madon.NewApp(godon.AppName, godon.WebSite, __SCOPES, __REDIRECT_URL, godon.InstanceName)

	if err != nil {
		return errors.Wrap(err, errMsg)
	}

	godon.client = client
	return nil
}

func (godon *Godon) restore() error {
	const errMsg = "godon.restore failed"

	client, err := madon.RestoreApp(godon.AppName, godon.InstanceName, godon.AppId, godon.AppSecret, godon.Token)

	if err != nil {
		return errors.Wrap(err, errMsg)
	}

	godon.client = client
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

func (godon *Godon) GetAppId() string {
	return godon.client.ID
}

func (godon *Godon) GetAppSecret() string {
	return godon.client.Secret
}

func (godon *Godon) PostStatus(toot Toot) error {
	const errMsg = "godon.PostStatus failed"
	// TODO Figure out who to send the status to.
	_, err := godon.client.PostStatus(toot.Status, toot.InReplyTo, toot.MediaIDs, toot.Sensitive, toot.SpoilerText, toot.Visibility)

	if err != nil {
		return errors.Wrap(err, errMsg)
	}

	return nil
}

func (godon *Godon) GetAccessToken() *madon.UserToken {
	if godon.client.UserToken == nil {
		return &madon.UserToken{}
	}

	return godon.client.UserToken
}

const __REDIRECT_URL = ""

var __SCOPES = []string{"read", "write", "follow"}

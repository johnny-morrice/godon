package godon

import ()

type Options struct {
}

type Godon struct {
	Options
}

func New(options Options) (*Godon, error) {
	panic("not implemented")

}

type AuthorizationMethod func(url string) (string, error)

func (godon *Godon) OAuthLogin(method AuthorizationMethod) error {
	panic("not implemented")
}

// TODO should take toot as parameter.
func (godon *Godon) PostStatus(statusText string) error {
	panic("not implemented")
}

const __REDIRECT_URL_VALUE = "urn:ietf:wg:oauth:2.0:oob"
const __REDIRECT_URL_KEY = "redirect_uri"
const __SCOPES_KEY = "scopes"
const __CLIENT_ID_KEY = "client_id"
const __RESPONSE_TYPE_KEY = "response_type"
const __RESPONSE_TYPE_VALUE = "code"

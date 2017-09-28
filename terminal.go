package godon

import (
	"fmt"

	ui "github.com/gizak/termui"
)

type Terminal struct {
	Client    *Client
	Body      *ui.Grid
	resources map[string]Resource
}

func (term Terminal) Render() Stopper {
	// FIXME returns error
	ui.Init()

	term.Body = ui.Body

	panic("not implemented")
}

func (term Terminal) AddResource(resource Resource) error {
	_, present := term.resources[resource.Name]

	if present {
		return fmt.Errorf("Duplicate resource name: %s", resource.Name)
	}

	term.resources[resource.Name] = resource
	return nil
}

type ShowWidget struct {
	Resource string
}

func (show ShowWidget) Execute(term Terminal) error {
	panic("not implemented")
}

type HideWidget struct {
	Resource string
}

func (hide HideWidget) Execute(term Terminal) error {
	panic("not implemented")
}

type Command interface {
	Execute(term Terminal) error
}

func (term Terminal) ExecuteCommand(cmd Command) error {
	panic("not implemented")
}

type Style uint16

const (
	Feed = iota
	Profile
	Input
)

type Resource struct {
	Name     string
	Style    Style
	Bufferer ui.GridBufferer
}

type Stopper struct {
	stopch chan<- struct{}
	errch  <-chan error
}

func (stopper Stopper) Stop() {
	stopper.stopch <- struct{}{}
}

func (stopper Stopper) Wait() error {
	return <-stopper.errch
}

type TimelineWidget struct {
	ui.Block
	Limit    int
	Local    bool
	Timeline string
}

package godon

import (
	"fmt"

	ui "github.com/gizak/termui"
)

type Terminal struct {
	Client    *Client
	resources map[string]Resource
	stopper   Stopper
}

func (term Terminal) Run() Stopper {
	term.stopper = makeStopper()

	err := ui.Init()

	if err != nil {
		term.fatal(err)
		return term.stopper
	}

	go func() {
		term.stopper.WaitForStop()
		ui.Close()
		ui.StopLoop()
	}()

	go ui.Loop()

	return term.stopper
}

func (term Terminal) close() {
	term.stopper.Stop()
}

func (term Terminal) fatal(err error) {
	term.stopper.Fail(err)
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
	stopch chan struct{}
	errch  chan error
}

func makeStopper() Stopper {
	return Stopper{
		stopch: make(chan struct{}),
		errch:  make(chan error, 1),
	}
}

func (stopper Stopper) Stop() {
	close(stopper.stopch)
	close(stopper.errch)
}

func (stopper Stopper) Fail(err error) {
	stopper.errch <- err
	stopper.Stop()
}

func (stopper Stopper) WaitForStop() {
	<-stopper.stopch
}

func (stopper Stopper) WaitForFailure() error {
	return <-stopper.errch
}

type TimelineWidget struct {
	ui.Block
	Limit    int
	Local    bool
	Timeline string
}

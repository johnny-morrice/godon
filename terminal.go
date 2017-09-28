package godon

import (
	"github.com/gizak/termui"
)

type Terminal struct {
	Client *Client
}

func (term Terminal) Render() Stopper {
	panic("not implemented")
}

func (term Terminal) AddResource(resource Resource) error {
	panic("not implemented")
}

type CommandOpCode uint16

const (
	ShowWidget = iota
	HideWidget
)

type Command struct {
	OpCode     CommandOpCode
	Parameters []string
}

func (term Terminal) ExecuteCommand(cmd Command) error {
	panic("not implemented")
}

type Resource struct {
	Name     string
	Bufferer termui.GridBufferer
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
	termui.Block
	Limit    int
	Local    bool
	Timeline string
}

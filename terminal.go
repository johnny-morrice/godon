package godon

type Terminal struct {
	Client *Client
}

func (term Terminal) Render() error {
	panic("not implemented")
}

type TimelineWidget struct {
	Limit    int
	Local    bool
	Timeline string
}

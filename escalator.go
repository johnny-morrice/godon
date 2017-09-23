package godon

// Escalator is a widget for displaying paginated data from a remote source.
//
// The Escalator manages three indices:
//
// 1. Logical index: the index of the data item in the remote source.
// 2. Physical index: the index of the data item in the list loaded so far by the client.
// 3. Visual index: the placement of the data item on screen, in the form of a nested widget.
type Escalator struct {
}

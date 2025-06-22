package model

// Stream represents either an input or output stream in a puzzle,
// including its type, name, position, and values.
type Stream struct {
	Type     StreamType
	Name     string
	Position uint8
	Values   []int16
}

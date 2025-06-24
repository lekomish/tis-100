package model

// Stream represents either an input or output stream in a puzzle,
// including its type, name, position, and values.
type Stream struct {
	Type     StreamType
	Name     string
	Position uint8
	Values   []int16
}

// Len returns the number of values in the stream.
// This indicates the length of the input or output data sequence.
func (s *Stream) Len() int {
	return len(s.Values)
}

package engine

import "github.com/lekomish/tis-100/internal/model"

// Output represents an output stream from a TIS-100 node.
// It stores all values sent to a specific output port during execution.
type Output struct {
	Index  uint8   // the index or ID of the output stream (e.g., 1)
	Values []int16 // the values written to this stream during simulation
}

// NewOutput initializes and returns a new Output instance
// with a given stream index and an empty value buffer.
func NewOutput(index uint8) *Output {
	return &Output{
		Index:  index,
		Values: make([]int16, 0),
	}
}

// AddValue appends a single value to the output stream.
func (o *Output) AddValue(value int16) {
	o.Values = append(o.Values, value)
}

// Len returns the number of values currently stored in the output stream.
func (o *Output) Len() int {
	return len(o.Values)
}

// At safely retrieves the value at the given index in the output stream.
// It returns false if the index is out of bounds.
func (o *Output) At(index int) (int16, bool) {
	if index < 0 || index >= o.Len() {
		return 0, false
	}
	return o.Values[index], true
}

// Clear resets the output stream by discarding all stored values,
// but retains the stream index.
func (o *Output) Clear() {
	o.Values = o.Values[:0]
}

// EqualToStream checks whether the Output matches the given Stream
// in position and values.
func (o *Output) EqualToStream(stream *model.Stream) bool {
	if stream == nil {
		return false
	}
	if o.Index != stream.Position {
		return false
	}
	if o.Len() != stream.Len() {
		return false
	}
	for i := 0; i < o.Len(); i++ {
		if o.Values[i] != stream.Values[i] {
			return false
		}
	}
	return true
}

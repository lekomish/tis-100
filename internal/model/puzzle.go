package model

// Puzzle defines a playable TIS-100 puzzle, including metadata, streams, and layout.
type Puzzle struct {
	Title       string
	Description []string
	Streams     []*Stream
	Layout      []NodeType
}

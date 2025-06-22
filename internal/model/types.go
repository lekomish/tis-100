package model

type (
	// StreamType defines whether a stream is an input or output.
	StreamType uint8
	// NodeType defines the type of a node used in the puzzle layout.
	NodeType uint8
)

const (
	// INPUT represents an input stream.
	INPUT StreamType = iota
	// OUTPUT represents an output stream.
	OUTPUT
)

const (
	// COMPUTE represents a functional node.
	COMPUTE NodeType = iota
	// DAMAGED represents a broken or unusable node.
	DAMAGED
)

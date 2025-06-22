// Package model defines the core data sturctures for TIS-100,
// including puzzles, code representations, streams, and node types.
package model

// Code represents a TIS-100 program with a title and a list of instructions per node.
type Code struct {
	Title string
	Nodes [][]string
}

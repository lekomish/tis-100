package engine

import "strings"

// InputCode represents raw assembly-like source code for a node,
// including the lines of code and any labels mapped to line numbers.
type InputCode struct {
	Lines  []string         // list of code lines, in source order
	Labels map[string]uint8 // mapping from label name to line index
}

// NewInputCode initializes and returns a pointer to an InputCode instance
// with preallocated slices and maps.
func NewInputCode() *InputCode {
	return &InputCode{
		Lines:  make([]string, 0),
		Labels: make(map[string]uint8),
	}
}

// AddLine appends a new line of code to the input.
// Empty or whitespace-only lines are ignored.
func (ic *InputCode) AddLine(line string) {
	line = strings.TrimSpace(line)
	if line == "" {
		return
	}
	ic.Lines = append(ic.Lines, line)
}

// AddLabel registers a label pointing to the given line index.
// If the label already exists, it ignored to avoid accidental overwrite.
func (ic *InputCode) AddLabel(name string, lineIndex uint8) {
	name = strings.TrimSpace(name)
	if name == "" {
		return
	}
	if _, exists := ic.Labels[name]; exists {
		return
	}
	ic.Labels[name] = lineIndex
}

// LineAt safely retrieves a line by index. Returns false if the index is out of bounds.
func (ic *InputCode) LineAt(index int) (string, bool) {
	if index < 0 || index >= len(ic.Lines) {
		return "", false
	}
	return ic.Lines[index], true
}

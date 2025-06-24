package engine

import "github.com/lekomish/tis-100/internal/model"

// advancePC increments the instruction pointer to move to the next instruction.
// It should be called after a successful instruction execution.
func (n *Node) advancePC() {
	n.InstructionPointer++
}

// jumpTo sets the instruction pointer to the specified position.
// If the position is out of bounds, it defaults to 0.
func (n *Node) jumpTo(pos int16) {
	if pos >= int16(len(n.Instructions)) || pos < 0 {
		pos = 0
	}
	n.InstructionPointer = uint8(pos)
}

// clampACC ensures that the ACC register stays within the defined bounds.
// If ACC exceeds `MaxACC` or falls below `MinACC`, it is clamped accordingly.
func (n *Node) clampACC() {
	if n.ACC > model.MaxACC {
		n.ACC = model.MaxACC
	}
	if n.ACC < model.MinACC {
		n.ACC = model.MinACC
	}
}

// resetPCIfOutOfBounds resets the instruction pointer to 0 if it points past the instruction list.
// This prevents out-of-bounds access and is typically used as a safety check before execution.
func (n *Node) resetPCIfOutOfBounds() {
	if n.InstructionPointer >= uint8(len(n.Instructions)) {
		n.InstructionPointer = 0
	}
}

// instruction returns the current instruction pointed to by the instruction pointer.
// If the pointer is out of bounds, it returns nil.
func (n *Node) instruction() *Instruction {
	if int(n.InstructionPointer) < len(n.Instructions) {
		return n.Instructions[n.InstructionPointer]
	}
	return nil
}

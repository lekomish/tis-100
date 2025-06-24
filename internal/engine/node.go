package engine

import "errors"

// Node represents a single compute node in the TIS-100 virtual machine.
// Each node has its own accumulator (ACC), backup register (BAK),
// instruction memory, and connections to neighboring nodes.
type Node struct {
	Index              uint8          // index of the node in the grid
	IsBlocked          bool           // whether the node is currently blocked (e.g., waiting for input/output)
	InstructionPointer uint8          // points to the current instruction being executed
	Instructions       []*Instruction // instruction memory for the node
	ACC                int16          // accumulator register
	BAK                int16          // backup register
	OutboundTarget     *Node          // node this node is currently writing to (if any)
	Last               *Node          // last node successfully communicated with
	OutboundValue      int16          // value being sent to `OutboundTarget`
	Ports              [4]*Node       // connections to neighboring nodes (UP, RIGHT, DOWN, LEFT)
	Output             *Output        // optional output collector for OUT instruction
}

// NewNode creates and returns a new Node initialized instruction memory and ports.
func NewNode() *Node {
	return &Node{
		Instructions: make([]*Instruction, 0),
		Ports:        [4]*Node{nil, nil, nil, nil},
	}
}

// Tick executes a single instruction cycle on the node.
// It fetches, decodes, and executes the current instruction.
// If the instruction is MOV, ADD, SUB, etc., it performs reads/writes as needed.
// If the instruction is a jump, it may update the instruction pointer.
// Returns an error if execution fails (e.g., invalid instruction or read/write error).
func (n *Node) Tick() error {
	// prevent execution if no instructions are loaded
	if len(n.Instructions) == 0 {
		return errors.New("no instructions to execute")
	}

	n.IsBlocked = true
	n.resetPCIfOutOfBounds()
	// fetch current instruction
	ins := n.instruction()

	switch ins.Op {
	case OpMov:
		// MOV SRC DEST - read and write a value
		val, blocked, err := n.read(ins.SrcType, ins.Src)
		if err != nil || blocked {
			return err
		}
		blocked, err = n.write(ins.Dest.Port, val)
		if err != nil || blocked {
			return err
		}
	case OpAdd:
		// ADD SRC - add value to ACC
		val, blocked, err := n.read(ins.SrcType, ins.Src)
		if err != nil || blocked {
			return err
		}
		n.ACC += val
		n.clampACC()
	case OpSub:
		// SUB SRC - subtract value from ACC
		val, blocked, err := n.read(ins.SrcType, ins.Src)
		if err != nil || blocked {
			return err
		}
		n.ACC -= val
		n.clampACC()
	case OpJmp:
		// JMP LABEL - unconditional jump
		n.jumpTo(ins.Src.Value)
		return nil
	case OpJro:
		// JRO OFFSET - relative jump
		n.jumpTo(int16(n.InstructionPointer) + ins.Src.Value)
		return nil
	case OpJez:
		// JEZ LABEL - jump if ACC == 0
		if n.ACC == 0 {
			n.jumpTo(ins.Src.Value)
			return nil
		}
	case OpJgz:
		// JGZ LABEL - jump if ACC > 0
		if n.ACC > 0 {
			n.jumpTo(ins.Src.Value)
			return nil
		}
	case OpJlz:
		// JLZ LABEL - jump if ACC < 0
		if n.ACC < 0 {
			n.jumpTo(ins.Src.Value)
			return nil
		}
	case OpJnz:
		// JNZ LABEL - jump if ACC != 0
		if n.ACC != 0 {
			n.jumpTo(ins.Src.Value)
			return nil
		}
	case OpSwp:
		// SWP - swap ACC and BAK
		n.ACC, n.BAK = n.BAK, n.ACC
	case OpSav:
		// SAV - copy ACC to BAK
		n.BAK = n.ACC
	case OpNeg:
		// NEG - negate ACC
		n.ACC *= -1
	case OpNop:
		// NOP - do nothing
	case OpOut:
		// OUT - write ACC value to output collector
		if n.Output != nil {
			n.Output.AddValue(n.ACC)
		}
	default:
		return errors.New("unknown operation")
	}

	// mark node as not blocked and move to the next instruction
	n.IsBlocked = false
	n.advancePC()
	return nil
}

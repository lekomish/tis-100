package engine

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// opCodeMap maps strings mnemonics to their corresponding OpCode values.
var opCodeMap = map[string]OpCode{
	"SUB": OpSub,
	"ADD": OpAdd,
	"JEZ": OpJez,
	"JMP": OpJmp,
	"JNZ": OpJnz,
	"JGZ": OpJgz,
	"JLZ": OpJlz,
	"JRO": OpJro,
	"SAV": OpSav,
	"SWP": OpSwp,
	"NOP": OpNop,
	"NEG": OpNeg,
	"OUT": OpOut,
}

// portNameMap maps strings representations of ports to their corresponding Port enum.
var portNameMap = map[string]Port{
	"UP":    PortUp,
	"DOWN":  PortDown,
	"LEFT":  PortLeft,
	"RIGHT": PortRight,
	"ACC":   PortAcc,
	"NIL":   PortNil,
	"ANY":   PortAny,
	"LAST":  PortLast,
}

// appendInstruction creates a new Instruction with the given OpCode
// and appends it to the node's instruction list.
func (n *Node) appendInstruction(op OpCode) *Instruction {
	ins := &Instruction{Op: op}
	n.Instructions = append(n.Instructions, ins)
	return ins
}

// compileCode compiles raw lines of TIS-100 code form InputCode into executable Instructions for the node.
// It also extracts and stores labels in the InputCode for jump resolution.
func (n *Node) compileCode(ic *InputCode) error {
	for i, line := range ic.Lines {
		if ind := strings.Index(line, ":"); ind != -1 {
			label := strings.TrimSpace(line[:ind])
			if label == "" {
				return fmt.Errorf("empty label in line %d", i)
			}
			ic.AddLabel(label, uint8(i))

			args := strings.TrimSpace(line[ind+1:])
			if len(args) == 0 {
				args = "NOP"
			}
			ic.Lines[i] = args
		}
	}

	for _, line := range ic.Lines {
		if err := n.parseInstruction(ic, line); err != nil {
			return err
		}
	}

	return nil
}

// parseInstruction determines the type of instruction and delegates to the appropriate parser.
func (n *Node) parseInstruction(ic *InputCode, line string) error {
	if len(line) < 3 {
		return errors.New("invalid line length")
	}

	mnemonic := strings.ToUpper(line[:3])
	if mnemonic == "MOV" {
		return n.parseMovInstruction(line)
	}

	op, ok := opCodeMap[mnemonic]
	if !ok {
		return fmt.Errorf("invalid instruction: %s", mnemonic)
	}

	switch op {
	case OpSub, OpAdd:
		return n.parseUnaryInstruction(line, op)
	case OpJez, OpJmp, OpJnz, OpJgz, OpJlz, OpJro:
		return n.parseJumpInstruction(ic, line, op)
	case OpSav, OpSwp, OpNop, OpNeg, OpOut:
		n.appendInstruction(op)
		return nil
	default:
		return fmt.Errorf("invalid instruction: %q", op)
	}
}

// parseMovInstruction parses a MOV instruction,
// which expects two operands (source and destination).
func (n *Node) parseMovInstruction(line string) error {
	if len(line) <= 3 {
		return errors.New("wrong mov instruction format")
	}

	args := strings.TrimSpace(line[4:])
	tokens := strings.FieldsFunc(args, func(r rune) bool {
		return r == ' ' || r == ','
	})
	if len(tokens) != 2 {
		return fmt.Errorf("MOV expects 2 arguments: %q", line)
	}

	ins := n.appendInstruction(OpMov)
	if err := parseOperation(tokens[0], &ins.SrcType, &ins.Src); err != nil {
		return err
	}
	if err := parseOperation(tokens[1], &ins.DestType, &ins.Dest); err != nil {
		return err
	}
	return nil
}

// parseUnaryInstruction parses single-operand instructions like ADD, SUB, etc.
func (n *Node) parseUnaryInstruction(line string, op OpCode) error {
	if len(line) < 4 {
		return fmt.Errorf("wrong instruction: %q", line)
	}
	args := strings.TrimSpace(line[4:])
	ins := n.appendInstruction(op)
	return parseOperation(args, &ins.SrcType, &ins.Src)
}

// parseJumpInstruction parses jump instructions using label references.
// It resolves labels from the InputCode's label map and sets the operand accordingly.
func (n *Node) parseJumpInstruction(ic *InputCode, line string, op OpCode) error {
	if len(line) < 4 {
		return fmt.Errorf("wrong instruction: %q", line)
	}
	label := strings.TrimSpace(line[4:])
	pos, ok := ic.Labels[label]
	if !ok {
		return fmt.Errorf("label not found: %s", label)
	}
	ins := n.appendInstruction(op)
	ins.SrcType = Immediate
	ins.Src.Value = int16(pos)
	return nil
}

// parseOperation parses an operand string and fills in the OperandType and Operand fields.
// Supports both port names and immediate numeric values.
func parseOperation(strOp string, opType *OperandType, op *Operand) error {
	if strOp == "" {
		return errors.New("no source was found")
	}

	*opType = PortRef
	if port, ok := portNameMap[strings.ToUpper(strOp)]; ok {
		op.Port = port
	} else {
		num, err := strconv.Atoi(strOp)
		if err != nil {
			return fmt.Errorf("invalid operand: %q", strOp)
		}
		*opType = Immediate
		op.Value = int16(num)
	}
	return nil
}

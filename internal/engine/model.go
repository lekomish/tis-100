package engine

type (
	// OpCode represents the operation to perform in an instruction.
	OpCode uint8

	// OperandType specifies whether an operand is a literal (Immediate)
	// or a port/register reference (PortRef).
	OperandType uint8

	// Port represents a communication direction or register within a node.
	Port uint8
)

// Operand defines a value or port used as input/output for an instruction.
// Interpretation depends on the associated OperandType.
type Operand struct {
	Value int16 // used when the operand is a literal (Immediate)
	Port  Port  // used when the operand is a port reference (PortRef)
}

// Instruction describes a complete TIS-100-sytle instruction,
// including the operation and its source and destination operands.
type Instruction struct {
	Op       OpCode      // the operation to execute (e.g., MOV, ADD)
	SrcType  OperandType // type of source operand (Immediate or PortRef)
	Src      Operand     // source operand value or port
	DestType OperandType // type of destination operand
	Dest     Operand     // destination operand value or port
}

// Supported OpCode values.
const (
	OpMov OpCode = iota // MOV: move value from source to destination
	OpSav               // SAV: save ACC to BAK
	OpNeg               // NEG: negate ACC
	OpAdd               // ADD: add source to ACC
	OpSub               // SUB: subtract source from ACC
	OpSwp               // SWP: swap ACC and BAK
	OpNop               // NOP: no operation

	// Conditional and relative jumps
	OpJmp // JMP: unconditional jump
	OpJro // JRO: relative jump by value
	OpJez // JEZ: jump if ACC == 0
	OpJnz // JNZ: jump if ACC != 0
	OpJlz // JLZ: jump if ACC < 0
	OpJgz // JGZ: jump if ACC > 0

	OpOut // OUT: write ACC to output
)

// OperandType values define how an operand is interpreted.
const (
	Immediate OperandType = iota // literal number (e.g., 5)
	PortRef                      // Port or register reference (e.g., ACC, UP)
)

// Port values used for inter-node communication and register access.
const (
	PortUp    Port = iota // northward port
	PortDown              // southward port
	PortLeft              // westward port
	PortRight             // eastward port
	PortAcc               // ACC register
	PortAny               // any connected port (e.g., for MOV ANY)
	PortLast              // the last successful communication port
	PortNil               // no connection/invalid
)

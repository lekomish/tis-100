package engine_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/lekomish/tis-100/internal/engine"
)

/* TESTS */

// --- NewNode ---
func TestNewNode(t *testing.T) {
	node := engine.NewNode()
	require.NotNil(t, node)
	require.Empty(t, node.Instructions)
	require.Equal(t, [4]*engine.Node{nil, nil, nil, nil}, node.Ports)
}

// --- TestTick ---
func TestNodeTickNoInstructions(t *testing.T) {
	node := engine.NewNode()
	err := node.Tick()
	require.ErrorContains(t, err, "no instructions to execute")
}

func TestNodeTickMovImmediateToAcc(t *testing.T) {
	node := engine.NewNode()
	node.Instructions = append(node.Instructions, movImmediateToAcc(42))

	err := node.Tick()
	require.NoError(t, err)
	require.Equal(t, int16(42), node.ACC)
	require.False(t, node.IsBlocked)
}

func TestNodeTickAdd(t *testing.T) {
	node := engine.NewNode()
	node.ACC = 10
	node.Instructions = append(node.Instructions, &engine.Instruction{
		Op:      engine.OpAdd,
		SrcType: engine.Immediate,
		Src:     engine.Operand{Value: 5},
	})

	err := node.Tick()
	require.NoError(t, err)
	require.Equal(t, int16(15), node.ACC)
}

func TestNodeTickSub(t *testing.T) {
	node := engine.NewNode()
	node.ACC = 20
	node.Instructions = append(node.Instructions, &engine.Instruction{
		Op:      engine.OpSub,
		SrcType: engine.Immediate,
		Src:     engine.Operand{Value: 7},
	})

	err := node.Tick()
	require.NoError(t, err)
	require.Equal(t, int16(13), node.ACC)
}

func TestNodeTickSavSwp(t *testing.T) {
	node := engine.NewNode()
	node.ACC = 9
	node.Instructions = append(node.Instructions, &engine.Instruction{Op: engine.OpSav})
	node.Instructions = append(node.Instructions, &engine.Instruction{Op: engine.OpSwp})

	_ = node.Tick()
	require.Equal(t, int16(9), node.BAK)

	_ = node.Tick()
	require.Equal(t, int16(9), node.ACC)
	require.Equal(t, int16(9), node.BAK)
}

func TestNodeTickNeg(t *testing.T) {
	node := engine.NewNode()
	node.ACC = -5
	node.Instructions = append(node.Instructions, &engine.Instruction{Op: engine.OpNeg})

	err := node.Tick()
	require.NoError(t, err)
	require.Equal(t, int16(5), node.ACC)
}

func TestNodeTickJumpInstructions(t *testing.T) {
	node := engine.NewNode()
	node.ACC = 5

	// 0: JMP 2
	// 1: NEG (should be skipped)
	// 2: NOP
	node.Instructions = []*engine.Instruction{
		{
			Op:      engine.OpJmp,
			SrcType: engine.Immediate,
			Src:     engine.Operand{Value: 2},
		},
		{Op: engine.OpNeg},
		{Op: engine.OpNop},
	}

	err := node.Tick()
	require.NoError(t, err)
	require.Equal(t, uint8(2), node.InstructionPointer)
	require.Equal(t, int16(5), node.ACC)
}

func TestNodeTickOut(t *testing.T) {
	node := engine.NewNode()
	node.ACC = 123
	node.Output = engine.NewOutput(0)

	node.Instructions = append(node.Instructions, &engine.Instruction{Op: engine.OpOut})

	err := node.Tick()
	require.NoError(t, err)
	require.Len(t, node.Output.Values, 1)
	require.Equal(t, int16(123), node.Output.Values[0])
}

func TestNodeTick_InvalidOp(t *testing.T) {
	node := engine.NewNode()
	node.Instructions = append(node.Instructions, &engine.Instruction{Op: engine.OpCode(255)})

	err := node.Tick()
	require.ErrorContains(t, err, "unknown operation")
}

// advancePC -> covered in previous tests
// jumpTo -> covered in previous tests
// clampACC -> covered in previous tests
// resetPCIfOutOfBounds -> covered in previous tests
// instruction -> covered in previous tests

/* UTILS */

// movImmediateToAcc creates a basic MOV instruction
func movImmediateToAcc(value int16) *engine.Instruction {
	return &engine.Instruction{
		Op:       engine.OpMov,
		SrcType:  engine.Immediate,
		Src:      engine.Operand{Value: value},
		DestType: engine.PortRef,
		Dest:     engine.Operand{Port: engine.PortAcc},
	}
}

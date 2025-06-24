package engine_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/lekomish/tis-100/internal/engine"
	"github.com/lekomish/tis-100/internal/model"
)

/* TESTS */

// --- NewEngine ---
func TestNewEngineValidSetup(t *testing.T) {
	code := &model.Code{
		Title: "TEST",
		Nodes: [][]string{
			{"MOV UP DOWN"}, // node 0
			{}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {},
		},
	}
	streams := []*model.Stream{
		{Type: model.INPUT, Position: 0, Values: []int16{1}},
		{Type: model.OUTPUT, Position: 0, Values: []int16{1}},
	}

	eng, err := engine.NewEngine(streams, code)
	require.NoError(t, err)
	require.NotNil(t, eng)
	require.Len(t, eng.Nodes, 12)
	require.Len(t, eng.Outputs, 1)
	require.Equal(t, uint8(0), eng.Outputs[0].Index)
}

func TestNewEngineInvalidCodeLength(t *testing.T) {
	code := &model.Code{
		Title: "SHORT",
		Nodes: [][]string{{"MOV UP DOWN"}},
	}
	strems := []*model.Stream{}

	eng, err := engine.NewEngine(strems, code)
	require.Nil(t, eng)
	require.ErrorContains(t, err, "wrong nodes number")
}

// --- Tick ---
func TestEngineTickAllNodesBlocked(t *testing.T) {
	code := &model.Code{
		Title: "TEST",
		Nodes: [][]string{
			{"MOV UP DOWN"},
			{},
			{},
			{},
			{},
			{},
			{},
			{},
			{},
			{},
			{},
			{},
		},
	}
	streams := []*model.Stream{}

	eng, err := engine.NewEngine(streams, code)
	require.NoError(t, err)

	blocked, tickErr := eng.Tick()
	require.NoError(t, tickErr)
	require.True(t, blocked)
}

func TestEngineTickWithInputOutput(t *testing.T) {
	code := &model.Code{
		Title: "SIMPLE-PIPE",
		Nodes: [][]string{
			{"MOV UP DOWN"},
			{},
			{},
			{},
			{"MOV UP DOWN"},
			{},
			{},
			{},
			{"MOV UP DOWN"},
			{},
			{},
			{},
		},
	}
	streams := []*model.Stream{
		{Type: model.INPUT, Position: 0, Values: []int16{42}},
		{Type: model.OUTPUT, Position: 0, Values: []int16{42}},
	}

	eng, err := engine.NewEngine(streams, code)
	require.NoError(t, err)

	for range 3 {
		_, err := eng.Tick()
		require.NoError(t, err)
	}

	require.Len(t, eng.Outputs[0].Values, 1)
	require.Equal(t, int16(42), eng.Outputs[0].Values[0])
}

// initStreams -> covered in previous tests
// loadInstructions -> covered in previous tests
// createEphemeralNode -> covered in previous tests
// createInputNode -> covered in previous tests
// createOutputNode -> covered in previous tests
// appendInstruction -> covered in previous tests
// compileCode -> covered in previous tests
// parseInstruction -> covered in previous tests
// parseMovInstruction -> covered in previous tests
// parseUnaryInstruction -> covered in previous tests
// parseJumpInstruction -> covered in previous tests
// parseOperation -> covered in previous tests
// read -> covered in previous tests
// write -> covered in previous tests
// getInputPort -> covered in previous tests
// getOutputPort -> covered in previous tests

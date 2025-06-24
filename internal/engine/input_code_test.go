package engine_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/lekomish/tis-100/internal/engine"
)

/* TESTS */

// --- NewInputCode ---
func TestNewInputCode(t *testing.T) {
	ic := engine.NewInputCode()
	require.NotNil(t, ic)
	require.Empty(t, ic.Lines)
	require.Empty(t, ic.Labels)
}

// --- AddLine ---
func TestAddLine(t *testing.T) {
	ic := engine.NewInputCode()

	ic.AddLine("MOV UP DOWN")
	require.Len(t, ic.Lines, 1)
	require.Equal(t, "MOV UP DOWN", ic.Lines[0])

	ic.AddLine("   ")
	ic.AddLine("")
	require.Len(t, ic.Lines, 1, "whitespace-only lines should be ignored")

	ic.AddLine(" ADD ACC 1 ")
	require.Len(t, ic.Lines, 2)
	require.Equal(t, "ADD ACC 1", ic.Lines[1])
}

// --- AddLabel ---
func TestAddLabel(t *testing.T) {
	ic := engine.NewInputCode()

	ic.AddLabel("start", 0)
	require.Len(t, ic.Labels, 1)
	require.Equal(t, uint8(0), ic.Labels["start"])

	ic.AddLabel("start", 1)
	require.Len(t, ic.Labels, 1, "duplicate label should be ignored")
	require.Equal(t, uint8(0), ic.Labels["start"], "duplicate label should be ignored")

	ic.AddLabel("   ", 2)
	ic.AddLabel("", 3)
	require.Len(t, ic.Labels, 1, "empty labels should be ignored")
}

// --- LineAt ---
func TestLineAt(t *testing.T) {
	ic := engine.NewInputCode()
	ic.AddLine("NOP")
	ic.AddLine("MOV 1 ACC")

	line, ok := ic.LineAt(0)
	require.True(t, ok)
	require.Equal(t, "NOP", line)

	line, ok = ic.LineAt(1)
	require.True(t, ok)
	require.Equal(t, "MOV 1 ACC", line)

	_, ok = ic.LineAt(2)
	require.False(t, ok)

	_, ok = ic.LineAt(-1)
	require.False(t, ok)
}

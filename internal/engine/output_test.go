package engine_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/lekomish/tis-100/internal/engine"
	"github.com/lekomish/tis-100/internal/model"
)

/* TESTS */

// --- NewOutput ---
func TestNewOutput(t *testing.T) {
	out := engine.NewOutput(2)
	require.Equal(t, uint8(2), out.Index)
	require.Equal(t, 0, out.Len())
}

// --- AddValue --- Len ---
func TestAddValueAndLen(t *testing.T) {
	out := engine.NewOutput(0)
	out.AddValue(5)
	out.AddValue(-3)

	require.Equal(t, 2, out.Len())
	require.Equal(t, int16(5), out.Values[0])
	require.Equal(t, int16(-3), out.Values[1])
}

// --- At ---
func TestAt(t *testing.T) {
	out := engine.NewOutput(0)
	out.AddValue(10)
	out.AddValue(20)

	v, ok := out.At(0)
	require.True(t, ok)
	require.Equal(t, int16(10), v)

	v, ok = out.At(1)
	require.True(t, ok)
	require.Equal(t, int16(20), v)

	_, ok = out.At(2)
	require.False(t, ok)

	_, ok = out.At(-1)
	require.False(t, ok)
}

// --- Clear ---
func TestClear(t *testing.T) {
	out := engine.NewOutput(1)
	out.AddValue(1)
	out.AddValue(2)

	require.Equal(t, 2, out.Len())
	out.Clear()
	require.Equal(t, 0, out.Len())
	require.Equal(t, uint8(1), out.Index)
}

// --- EqualToStream ---
func TestEqualToStream(t *testing.T) {
	out := engine.NewOutput(1)
	out.AddValue(100)
	out.AddValue(-50)

	stream := &model.Stream{
		Type:     model.OUTPUT,
		Name:     "TestOut",
		Position: 1,
		Values:   []int16{100, -50},
	}

	require.True(t, out.EqualToStream(stream))
}

func TestEqualToStreamWithDifferentIndex(t *testing.T) {
	out := engine.NewOutput(1)
	out.AddValue(100)

	stream := &model.Stream{
		Type:     model.OUTPUT,
		Name:     "WrongPos",
		Position: 2,
		Values:   []int16{100},
	}

	require.False(t, out.EqualToStream(stream))
}

func TestEqualToStreamWithDifferentValues(t *testing.T) {
	out := engine.NewOutput(1)
	out.AddValue(100)

	stream := &model.Stream{
		Type:     model.OUTPUT,
		Name:     "WrongVal",
		Position: 1,
		Values:   []int16{200},
	}

	require.False(t, out.EqualToStream(stream))
}

func TestEqualToStreamWithDifferentLengths(t *testing.T) {
	out := engine.NewOutput(1)
	out.AddValue(1)
	out.AddValue(2)

	stream := &model.Stream{
		Type:     model.OUTPUT,
		Name:     "Shorter",
		Position: 1,
		Values:   []int16{1},
	}

	require.False(t, out.EqualToStream(stream))
}

func TestEqualToStreamWithNilStream(t *testing.T) {
	out := engine.NewOutput(1)
	require.False(t, out.EqualToStream(nil))
}

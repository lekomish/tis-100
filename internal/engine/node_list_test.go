package engine_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/lekomish/tis-100/internal/engine"
)

/* TESTS */

// --- Append ---
func TestAppendToNilList(t *testing.T) {
	n := newTestNode(1)
	var list *engine.NodeList

	list = list.Append(n)

	require.NotNil(t, list)
	require.Equal(t, uint8(1), list.Node.Index)
	require.Nil(t, list.Next)
}

func TestAppendToNonEmptyList(t *testing.T) {
	n1 := newTestNode(1)
	n2 := newTestNode(2)

	list := (&engine.NodeList{Node: n1}).Append(n2)

	require.Equal(t, uint8(1), list.Node.Index)
	require.NotNil(t, list.Next)
	require.Equal(t, uint8(2), list.Next.Node.Index)
	require.Nil(t, list.Next.Next)
}

// --- Prepend ---
func TestPrependToNilList(t *testing.T) {
	n := newTestNode(3)
	var list *engine.NodeList

	list = list.Prepend(n)

	require.NotNil(t, list)
	require.Equal(t, uint8(3), list.Node.Index)
	require.Nil(t, list.Next)
}

func TestPrependToNonEmptyList(t *testing.T) {
	n1 := newTestNode(1)
	n2 := newTestNode(2)

	list := &engine.NodeList{Node: n1}
	list = list.Prepend(n2)

	require.Equal(t, uint8(2), list.Node.Index)
	require.NotNil(t, list.Next)
	require.Equal(t, uint8(1), list.Next.Node.Index)
	require.Nil(t, list.Next.Next)
}

/* UTILS */
func newTestNode(index uint8) *engine.Node {
	n := engine.NewNode()
	n.Index = index
	return n
}

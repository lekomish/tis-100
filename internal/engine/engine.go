// Package engine implements the core execution logic for simulationg TIS-100-like nodes.
//
// The engine operates a grid of interconnected compute nodes arranged in a 4x3 layout,
// where each node executes a small set of assembly-like instructions. In addition to
// physical compute nodes, the engine also manages ephemeral input and output nodes
// that injects and capture values via streams.
//
// Each node operates independently per tick, using message-passing through directional
// ports (UP, DOWN, LEFT, RIGHT), shared accumulator (ACC), and auxiliary register (BAK).
//
// The engine tracks all active nodes, executes them in order each cycle, and reports
// when all are blocked (indicating program stalling or completion).
package engine

import (
	"errors"
	"strings"

	"github.com/lekomish/tis-100/internal/model"
)

const (
	rows           = 3           // number of rows in the node grid
	cols           = 4           // number of columns in the node grid
	nodesTotal     = rows * cols // total number of nodes in the grid
	positionOffset = 4           // offset for the node positioning below ephemeral input nodes
	outputOffset   = 16          // offset for the ephemeral output node positioning below actual compute node
)

// Engine represents the execution engine simulating the TIS-100 node grid,
// including runtime state, connections, and stream I/O.
type Engine struct {
	Nodes       []*Node   // all physical nodes in the gird
	NodeList    *NodeList // linked list of ephemeral input/output nodes
	ActiveNodes *NodeList // linked list of nodes that are acitve each tick
	Outputs     []*Output // output values produced by output nodes
}

// NewEngine initializes an `Engine` with the provided input/output streams and code.
// It creates and connects nodes in a 4x3 grid and loads the program into them.
func NewEngine(streams []*model.Stream, code *model.Code) (*Engine, error) {
	nodes := make([]*Node, 0, model.NodesNumber)
	for i := range model.NodesNumber {
		n := NewNode()
		n.Index = uint8(i + positionOffset)
		nodes = append(nodes, n)
	}

	e := &Engine{Nodes: nodes, Outputs: make([]*Output, 0)}
	// set up directional connections between adjacent nodes in the grid
	for i, n := range e.Nodes {
		row := i / cols
		col := i % cols
		if row < rows-1 {
			n.Ports[PortDown] = e.Nodes[i+cols]
		}
		if row > 0 {
			n.Ports[PortUp] = e.Nodes[i-cols]
		}
		if col < cols-1 {
			n.Ports[PortRight] = e.Nodes[i+1]
		}
		if col > 0 {
			n.Ports[PortLeft] = e.Nodes[i-1]
		}
	}

	if err := e.loadInstructions(code); err != nil {
		return nil, err
	}
	if err := e.initStreams(streams); err != nil {
		return nil, err
	}

	return e, nil
}

// Tick executes one cycle of the engine by ticking all active nodes.
// Returns true if all active nodes are blocked (i.e., no further progress is possible).
func (e *Engine) Tick() (bool, error) {
	allBlocked := true
	for list := e.ActiveNodes; list != nil; list = list.Next {
		if err := list.Node.Tick(); err != nil {
			return false, err
		}
		allBlocked = allBlocked && list.Node.IsBlocked
	}
	return allBlocked, nil
}

// initStreams initializes the stream nodes (input/output) and appends them to the active list.
func (e *Engine) initStreams(streams []*model.Stream) error {
	for _, stream := range streams {
		switch stream.Type {
		case model.INPUT:
			n := e.createInputNode(stream)
			e.ActiveNodes = e.ActiveNodes.Prepend(n)
		case model.OUTPUT:
			n := e.createOutputNode(stream)
			e.ActiveNodes = e.ActiveNodes.Append(n)
		default:
			return errors.New("unknown stream type")
		}
	}
	return nil
}

// loadInstructions parses and loads the given code into the corresponding physical nodes.
func (e *Engine) loadInstructions(code *model.Code) error {
	if len(code.Nodes) != model.NodesNumber {
		return errors.New("wrong nodes number")
	}

	allInput := make([]*InputCode, model.NodesNumber)
	for i := range model.NodesNumber {
		allInput[i] = NewInputCode()
	}

	// format and parse each line of code into uppercased instrucionts
	for i, n := range code.Nodes {
		for _, line := range n {
			formatted := strings.ToUpper(strings.TrimSpace(line))
			allInput[i].AddLine(formatted)
		}
	}

	// compile instructions for each node and mark them as active if needed
	for i, n := range e.Nodes {
		if err := n.compileCode(allInput[i]); err != nil {
			return err
		}
		if len(n.Instructions) > 0 {
			e.ActiveNodes = e.ActiveNodes.Append(n)
		}
	}

	return nil
}

// createEphemeralNode creates a temporary node (used for I/O) and adds it to the engine's node list.
func (e *Engine) createEphemeralNode() *Node {
	n := NewNode()
	e.NodeList = e.NodeList.Append(n)
	return n
}

// createInputNode constructs an input node that injects values from the stream into the node below.
func (e *Engine) createInputNode(stream *model.Stream) *Node {
	inputNode := e.createEphemeralNode()
	inputNode.Index = stream.Position
	belowNode := e.Nodes[stream.Position]

	// connect ports
	inputNode.Ports[PortDown] = belowNode
	belowNode.Ports[PortUp] = inputNode

	// load MOV instructions for each stream value
	for _, value := range stream.Values {
		ins := inputNode.appendInstruction(OpMov)
		ins.SrcType = Immediate
		ins.Src.Value = value
		ins.DestType = PortRef
		ins.Dest.Port = PortDown
	}

	// add a JRO 0 to keep the node busy and cycling
	ins := inputNode.appendInstruction(OpJro)
	ins.SrcType = Immediate
	ins.Src.Value = 0

	return inputNode
}

// createOutputNode constructs a node that reads from the node above and stores values in an `Output`.
func (e *Engine) createOutputNode(stream *model.Stream) *Node {
	outputNode := e.createEphemeralNode()
	outputNode.Index = stream.Position + outputOffset
	aboveNode := e.Nodes[stream.Position+outputOffset-2*positionOffset]

	// connect ports
	outputNode.Ports[PortUp] = aboveNode
	aboveNode.Ports[PortDown] = outputNode

	// instruction to pull value from above and store it in ACC
	ins := outputNode.appendInstruction(OpMov)
	ins.SrcType = PortRef
	ins.Src.Port = PortUp
	ins.DestType = PortRef
	ins.Dest.Port = PortAcc
	// instruction to output the ACC value
	outputNode.appendInstruction(OpOut)

	// bind output buffer to output node
	e.Outputs = append(e.Outputs, NewOutput(stream.Position))
	outputNode.Output = e.Outputs[len(e.Outputs)-1]

	return outputNode
}

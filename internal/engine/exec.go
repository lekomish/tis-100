package engine

import (
	"errors"
	"fmt"
)

// read attempts to read a value from the specified operand.
// It returns the value, a `blocked` flag (true if the read couldn't be completed),
// and an error if the read is malformed.
//
// Behavior:
// - Immediate: returns the constant value.
// - Port:
//   - NIL: returns 0.
//   - ACC: returns the node's ACC value.
//   - ANY: probes ports in order (Left, Right, Up, Down) and reads if one is writing to this node.
//   - LAST: returns 0 if a previous input node was used, but does not block.
//   - Otherwise: reads from the designated direction if the peer is targeting this node.
func (n *Node) read(opType OperandType, op Operand) (int16, bool, error) {
	if n.OutboundTarget != nil {
		// if this node is trying to send a value, it can't receive
		return 0, false, nil
	}
	if opType == Immediate {
		// literal value
		return op.Value, false, nil
	}

	switch op.Port {
	case PortNil:
		// NIL port reads return zero
		return 0, false, nil
	case PortAcc:
		// return from accumulator
		return n.ACC, false, nil
	case PortUp, PortRight, PortDown, PortLeft, PortAny, PortLast:
		readFrom := n.getInputPort(op.Port)
		if readFrom == nil {
			// no connection or no node sending data
			return 0, true, nil
		} else if readFrom.OutboundTarget == n {
			// successful read
			val := readFrom.OutboundValue

			// reset the sender
			readFrom.OutboundValue = 0
			readFrom.OutboundTarget = nil
			readFrom.advancePC()

			if op.Port == PortAny {
				n.Last = readFrom
			}
			return val, false, nil
		} else if op.Port == PortLast {
			// no data available from previous connection
			return 0, false, nil
		} else {
			// data not ready
			return 0, true, nil
		}
	default:
		return 0, false, fmt.Errorf("invalid port: %v", op.Port)
	}
}

// write attempts to send a value to the specified port.
// Returns a `blocked` flag (true if the send couldn't be completed) and an error if invalid.
//
// Behavior:
// - ACC: stores the value in the accumulator.
// - ANY: sends to the first compatible output node.
// - LAST: reuses the last successful destination.
// - Specific direction: sends to the connected node, if available and not already sending.
func (n *Node) write(port Port, value int16) (bool, error) {
	switch port {
	case PortAcc:
		n.ACC = value
		return false, nil
	case PortUp, PortRight, PortDown, PortLeft, PortAny, PortLast:
		dest := n.getOutputPort(port)
		if dest != nil && n.OutboundTarget == nil {
			n.OutboundTarget = dest
			n.OutboundValue = value
			if port == PortAny {
				n.Last = dest
			}
		}
		// couldn't find destination or already outputting
		return true, nil
	case PortNil:
		return false, errors.New("unable to write")
	default:
		return false, errors.New("nowhere to write")
	}
}

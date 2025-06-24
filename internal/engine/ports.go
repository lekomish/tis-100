package engine

// portProbeOrder defines the order in which ports are probed
// when using the PortAny directive. This affects how input and output
// behaviour resolves in ambiguous or multi-port connections.
var portProbeOrder = []Port{PortLeft, PortRight, PortUp, PortDown}

// getInputPort returns the appropriate input Node for the given port.
// Special handling:
// - `PortAny`: iterates over ports in `portProbeOrder` to find the first connected node sending to this node.
// - `PortLast`: returns the last node that successfully sent data to this one.
// - `Default`: return the node connected to the specified direction.
func (n *Node) getInputPort(port Port) *Node {
	switch port {
	case PortAny:
		for _, p := range portProbeOrder {
			node := n.Ports[p]
			if node != nil && node.OutboundTarget == n {
				return node
			}
		}
		return nil
	case PortLast:
		return n.Last
	default:
		return n.Ports[port]
	}
}

// getOutputPort returns the destination Node for the given output port.
// Special handling:
// - `PortAny`: searches for a node whose current instruction is MOV and whose source refers to this node.
// - `PortLast`: returns the last output node used.
// - `Default`: returns the node connected to the specified direction.
func (n *Node) getOutputPort(port Port) *Node {
	if port == PortAny {
		for _, p := range portProbeOrder {
			node := n.Ports[p]
			if node == nil {
				continue
			}

			ins := node.instruction()
			if ins == nil || ins.Op != OpMov || ins.SrcType != PortRef {
				continue
			}
			if ins.Src.Port == PortAny || node.Ports[ins.Src.Port] == n {
				return node
			}
		}
		return nil
	} else if port == PortLast {
		return n.Last
	} else {
		return n.Ports[port]
	}
}

package engine

// NodeList is a singly linked list of Node pointers.
type NodeList struct {
	Node *Node
	Next *NodeList
}

// Append adds a new node to the end of the list.
// Returns the head of the list (which may be the new node if list was nil).
func (list *NodeList) Append(n *Node) *NodeList {
	tail := &NodeList{Node: n}

	if list == nil {
		return tail
	}

	head := list
	for head.Next != nil {
		head = head.Next
	}
	head.Next = tail

	return list
}

// Prepend adds a new node to the beginning of the list.
// Returns the new head of the list.
func (list *NodeList) Prepend(n *Node) *NodeList {
	return &NodeList{
		Node: n,
		Next: list,
	}
}

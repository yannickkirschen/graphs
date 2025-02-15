package graphs

import "fmt"

type Node[I, P comparable] struct {
	id          I
	connections map[P][]P
}

func NewNode[I, P comparable](id I) *Node[I, P] {
	return &Node[I, P]{id, map[P][]P{}}
}

func (node *Node[I, P]) Id() I {
	return node.id
}

func (node *Node[I, P]) Connect(from, to P) {
	connection, ok := node.connections[from]
	if !ok {
		connection = []P{}
	}

	node.connections[from] = append(connection, to)
}

func (node *Node[I, P]) ConnectBi(from, to P) {
	node.Connect(from, to)
	node.Connect(to, from)
}

func (node *Node[I, P]) Next(port P) []P {
	return node.connections[port]
}

func (node *Node[I, P]) String() string {
	return fmt.Sprintf("Node<%v>", node.id)
}

func (node *Node[I, P]) Equals(o *Node[I, P]) bool {
	return node.id == o.id
}

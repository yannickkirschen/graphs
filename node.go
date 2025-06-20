package graphs

import "fmt"

type Node[O, P comparable] struct {
	id          O
	connections map[P][]P
}

func NewNode[O, P comparable](id O) *Node[O, P] {
	return &Node[O, P]{id, map[P][]P{}}
}

func (node *Node[O, P]) Id() O {
	return node.id
}

func (node *Node[O, P]) Connect(from, to P) {
	connection, ok := node.connections[from]
	if !ok {
		connection = []P{}
	}

	node.connections[from] = append(connection, to)
}

func (node *Node[O, P]) ConnectBi(from, to P) {
	node.Connect(from, to)
	node.Connect(to, from)
}

func (node *Node[O, P]) Next(port P) []P {
	return node.connections[port]
}

func (node *Node[O, P]) String() string {
	return fmt.Sprintf("Node<%v>", node.id)
}

func (node *Node[O, P]) Equals(o *Node[O, P]) bool {
	return node.id == o.id
}

package graphs

import "fmt"

type Connection[O, P comparable] struct {
	FromNode *Node[O, P]
	FromPort P
	ToNode   *Node[O, P]
	ToPort   P
}

func (connection *Connection[O, P]) IsSelf() bool {
	return connection.FromNode.Equals(connection.ToNode) && connection.FromPort == connection.ToPort
}

func (connection *Connection[O, P]) EqualNodes() bool {
	return connection.FromNode.Equals(connection.ToNode)
}

func (connection *Connection[O, P]) String() string {
	return fmt.Sprintf("%v.%v -> %v.%v", connection.FromNode.id, connection.FromPort, connection.ToNode.id, connection.ToPort)
}

func NewConnection[O, P comparable](fromNode *Node[O, P], fromPort P, toNode *Node[O, P], toPort P) *Connection[O, P] {
	return &Connection[O, P]{fromNode, fromPort, toNode, toPort}
}

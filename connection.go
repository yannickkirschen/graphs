package graphs

import "fmt"

type Connection[I, P comparable] struct {
	FromNode *Node[I, P]
	FromPort P
	ToNode   *Node[I, P]
	ToPort   P
}

func (connection *Connection[I, P]) IsSelf() bool {
	return connection.FromNode.Equals(connection.ToNode) && connection.FromPort == connection.ToPort
}

func (connection *Connection[I, P]) EqualNodes() bool {
	return connection.FromNode.Equals(connection.ToNode)
}

func (connection *Connection[I, P]) String() string {
	return fmt.Sprintf("%v.%v -> %v.%v", connection.FromNode.id, connection.FromPort, connection.ToNode.id, connection.ToPort)
}

func NewConnection[I, P comparable](fromNode *Node[I, P], fromPort P, toNode *Node[I, P], toPort P) *Connection[I, P] {
	return &Connection[I, P]{fromNode, fromPort, toNode, toPort}
}

package graphs

import (
	"fmt"
	"slices"
)

type Graph[I, P comparable] struct {
	nodes       map[I]*Node[I, P]
	connections []*Connection[I, P]
}

func NewGraph[I, P comparable]() *Graph[I, P] {
	return &Graph[I, P]{map[I]*Node[I, P]{}, []*Connection[I, P]{}}
}

func (graph *Graph[I, P]) AddNode(node *Node[I, P]) {
	graph.nodes[node.id] = node
}

func (graph *Graph[I, P]) AddConnection(connection *Connection[I, P]) error {
	if slices.Contains(graph.connections, connection) {
		return fmt.Errorf("graph: connection %s already exists and cannot be re-added", connection)
	}

	graph.connections = append(graph.connections, connection)
	return nil
}

func (graph *Graph[I, P]) Connect(fromNode *Node[I, P], fromPort P, toNode *Node[I, P], toPort P) error {
	return graph.AddConnection(NewConnection(fromNode, fromPort, toNode, toPort))
}

func (graph *Graph[I, P]) ConnectBi(fromNode *Node[I, P], fromPort P, toNode *Node[I, P], toPort P) error {
	if err := graph.AddConnection(NewConnection(fromNode, fromPort, toNode, toPort)); err != nil {
		return err
	}

	return graph.AddConnection(NewConnection(toNode, toPort, fromNode, fromPort))
}

func (graph *Graph[I, P]) ConnectRef(fromRef I, fromPort P, toRef I, toPort P) error {
	fromNode, ok := graph.nodes[fromRef]
	if !ok {
		return fmt.Errorf("graph: from ref %v not found and thus cannot be used for connection", fromRef)
	}

	toNode, ok := graph.nodes[toRef]
	if !ok {
		return fmt.Errorf("graph: to ref %v not found and thus cannot be used for connection", toRef)
	}

	return graph.Connect(fromNode, fromPort, toNode, toPort)
}

func (graph *Graph[I, P]) ConnectRefBi(fromRef I, fromPort P, toRef I, toPort P) error {
	if err := graph.ConnectRef(fromRef, fromPort, toRef, toPort); err != nil {
		return err
	}

	return graph.ConnectRef(toRef, toPort, fromRef, fromPort)
}

func (graph *Graph[I, P]) FindConnection(node *Node[I, P], port P) (*Connection[I, P], bool) {
	for _, connection := range graph.connections {
		if connection.FromNode.Equals(node) && connection.FromPort == port {
			return connection, true
		}
	}

	return nil, false
}

func (graph *Graph[I, P]) Find(fromNode *Node[I, P], fromPort P, toNode *Node[I, P], toPort P) [][]*Triple[*P, *Node[I, P], *P] {
	data := &dFSData[I, P]{[]*Node[I, P]{}, []*Triple[*P, *Node[I, P], *P]{}, [][]*Triple[*P, *Node[I, P], *P]{}}
	graph.dfsFind(data, nil, NewConnection(fromNode, fromPort, toNode, toPort))
	return data.Paths
}

func (graph *Graph[I, P]) FindRef(fromRef I, fromPort P, toRef I, toPort P) ([][]*Triple[*P, *Node[I, P], *P], error) {
	fromNode, ok := graph.nodes[fromRef]
	if !ok {
		return nil, fmt.Errorf("graph: from ref %v not found and thus cannot find paths", fromRef)
	}

	toNode, ok := graph.nodes[toRef]
	if !ok {
		return nil, fmt.Errorf("graph: to ref %v not found and thus cannot find paths", toRef)
	}

	return graph.Find(fromNode, fromPort, toNode, toPort), nil
}

func (graph *Graph[I, P]) dfsFind(data *dFSData[I, P], _ *Connection[I, P], current *Connection[I, P]) {
	if slices.Contains(data.Visited, current.FromNode) {
		return
	}

	var entry *P
	if len(data.CurrentPath)-1 >= 0 {
		previous, ok := graph.FindConnection(data.CurrentPath[len(data.CurrentPath)-1].Middle, *data.CurrentPath[len(data.CurrentPath)-1].Right)

		if ok {
			entry = &previous.ToPort
		}
	}

	data.Visited = append(data.Visited, current.FromNode)
	data.CurrentPath = append(data.CurrentPath, NewTriple(entry, current.FromNode, &current.FromPort))

	if current.IsSelf() {
		graph.dfsHandleFound(data)
		return
	}

	next, ok := graph.FindConnection(current.FromNode, current.FromPort)
	if !ok {
		data.CurrentPath = data.CurrentPath[:len(data.CurrentPath)-1]
		return
	}

	for _, port := range next.ToNode.Next(next.ToPort) {
		// This happens, when the next port belongs to the same node as the
		// current one.
		if current.EqualNodes() && port == current.ToPort {
			data.CurrentPath[len(data.CurrentPath)-1].Right = nil // Clear last port
			graph.dfsHandleFound(data)
			return
		}

		prev := NewConnection(current.FromNode, current.FromPort, next.FromNode, next.ToPort)
		new := NewConnection(next.ToNode, port, current.ToNode, current.ToPort)
		graph.dfsFind(data, prev, new)
	}

	data.CurrentPath = data.CurrentPath[:len(data.CurrentPath)-1]

	idx := slices.Index(data.Visited, current.FromNode)
	data.Visited = slices.Delete(data.Visited, idx, idx+1)
}

func (graph *Graph[I, P]) dfsHandleFound(data *dFSData[I, P]) {
	var currentPath []*Triple[*P, *Node[I, P], *P]
	currentPath = append(currentPath, data.CurrentPath...)

	data.Paths = append(data.Paths, currentPath)
	data.Visited = data.Visited[:len(data.Visited)-1]
	data.CurrentPath = data.CurrentPath[:len(data.CurrentPath)-1]
}

type dFSData[I comparable, P comparable] struct {
	Visited     []*Node[I, P]
	CurrentPath []*Triple[*P, *Node[I, P], *P]
	Paths       [][]*Triple[*P, *Node[I, P], *P]
}

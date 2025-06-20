package graphs

import (
	"fmt"
	"slices"

	"github.com/moznion/go-optional"
)

type PathSegment[O, P comparable] Triple[optional.Option[P], *Node[O, P], optional.Option[P]]

func (seg *PathSegment[O, P]) String() string {
	return fmt.Sprintf("<%v, %v, %v>", seg.Left, seg.Middle, seg.Right)
}

type Graph[O, P comparable] struct {
	nodes       map[O]*Node[O, P]
	connections []*Connection[O, P]
}

func NewGraph[O, P comparable]() *Graph[O, P] {
	return &Graph[O, P]{map[O]*Node[O, P]{}, []*Connection[O, P]{}}
}

func (graph *Graph[O, P]) AddNode(node *Node[O, P]) {
	graph.nodes[node.id] = node
}

func (graph *Graph[O, P]) AddConnection(connection *Connection[O, P]) error {
	if slices.Contains(graph.connections, connection) {
		return fmt.Errorf("graph: connection %s already exists and cannot be re-added", connection)
	}

	graph.connections = append(graph.connections, connection)
	return nil
}

func (graph *Graph[O, P]) Connect(fromNode *Node[O, P], fromPort P, toNode *Node[O, P], toPort P) error {
	return graph.AddConnection(NewConnection(fromNode, fromPort, toNode, toPort))
}

func (graph *Graph[O, P]) ConnectBi(fromNode *Node[O, P], fromPort P, toNode *Node[O, P], toPort P) error {
	if err := graph.AddConnection(NewConnection(fromNode, fromPort, toNode, toPort)); err != nil {
		return err
	}

	return graph.AddConnection(NewConnection(toNode, toPort, fromNode, fromPort))
}

func (graph *Graph[O, P]) ConnectRef(fromRef O, fromPort P, toRef O, toPort P) error {
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

func (graph *Graph[O, P]) ConnectRefBi(fromRef O, fromPort P, toRef O, toPort P) error {
	if err := graph.ConnectRef(fromRef, fromPort, toRef, toPort); err != nil {
		return err
	}

	return graph.ConnectRef(toRef, toPort, fromRef, fromPort)
}

func (graph *Graph[O, P]) FindConnection(node *Node[O, P], port P) (*Connection[O, P], bool) {
	for _, connection := range graph.connections {
		if connection.FromNode.Equals(node) && connection.FromPort == port {
			return connection, true
		}
	}

	return nil, false
}

func (graph *Graph[O, P]) Find(fromNode *Node[O, P], fromPort P, toNode *Node[O, P], toPort P) [][]*PathSegment[O, P] {
	data := &dFSData[O, P]{[]*Node[O, P]{}, []*PathSegment[O, P]{}, [][]*PathSegment[O, P]{}}
	graph.dfsFind(data, NewConnection(fromNode, fromPort, toNode, toPort))
	return data.Paths
}

func (graph *Graph[O, P]) FindRef(fromRef O, fromPort P, toRef O, toPort P) ([][]*PathSegment[O, P], error) {
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

func (graph *Graph[O, P]) dfsFind(data *dFSData[O, P], current *Connection[O, P]) {
	if slices.Contains(data.Visited, current.FromNode) {
		return
	}

	var entry P
	if len(data.CurrentPath)-1 >= 0 {
		previous, ok := graph.FindConnection(data.CurrentPath[len(data.CurrentPath)-1].Middle, data.CurrentPath[len(data.CurrentPath)-1].Right.Unwrap())

		if ok {
			entry = previous.ToPort
		}
	}

	data.Visited = append(data.Visited, current.FromNode)
	data.CurrentPath = append(data.CurrentPath, &PathSegment[O, P]{optional.Some(entry), current.FromNode, optional.Some(current.FromPort)})

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
			data.CurrentPath[len(data.CurrentPath)-1].Right = optional.None[P]() // Clear last port
			graph.dfsHandleFound(data)
			return
		}

		new := NewConnection(next.ToNode, port, current.ToNode, current.ToPort)
		graph.dfsFind(data, new)
	}

	data.CurrentPath = data.CurrentPath[:len(data.CurrentPath)-1]

	idx := slices.Index(data.Visited, current.FromNode)
	data.Visited = slices.Delete(data.Visited, idx, idx+1)
}

func (graph *Graph[O, P]) dfsHandleFound(data *dFSData[O, P]) {
	var currentPath []*PathSegment[O, P]
	currentPath = append(currentPath, data.CurrentPath...)

	data.Paths = append(data.Paths, currentPath)
	data.Visited = data.Visited[:len(data.Visited)-1]
	data.CurrentPath = data.CurrentPath[:len(data.CurrentPath)-1]
}

type dFSData[O comparable, P comparable] struct {
	Visited     []*Node[O, P]
	CurrentPath []*PathSegment[O, P]
	Paths       [][]*PathSegment[O, P]
}

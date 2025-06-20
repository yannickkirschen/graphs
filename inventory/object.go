package inventory

import (
	"github.com/moznion/go-optional"
	"github.com/yannickkirschen/graphs"
)

type Object[O, C, P comparable] struct {
	Id    O
	Label string
	Class *Class[C, P]
	Spec  optional.Option[any]
}

func NewObject[O, C, P comparable](id O, label string) *Object[O, C, P] {
	return &Object[O, C, P]{
		id,
		label,
		nil,
		optional.None[any](),
	}
}

func (object *Object[O, C, P]) ToGraphNode() *graphs.Node[O, P] {
	node := graphs.NewNode[O, P](object.Id)

	for _, connection := range object.Class.Connections {
		if connection.Bidirectional {
			node.ConnectBi(connection.From, connection.To)
		} else {
			node.Connect(connection.From, connection.To)
		}
	}

	return node
}

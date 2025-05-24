package inventory

import (
	"github.com/moznion/go-optional"
	"github.com/yannickkirschen/graphs"
)

type Object struct {
	Id    Id
	Label string
	Class *Class
	Spec  optional.Option[any]
}

func NewObject(id Id, label string) *Object {
	return &Object{
		id,
		label,
		nil,
		optional.None[any](),
	}
}

func (object *Object) ToGraphNode() *graphs.Node[Id, Id] {
	node := graphs.NewNode[Id, Id](object.Id)

	for _, connection := range object.Class.Connections {
		if connection.Bidirectional {
			node.ConnectBi(connection.From, connection.To)
		} else {
			node.Connect(connection.From, connection.To)
		}
	}

	return node
}

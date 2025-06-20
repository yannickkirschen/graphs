package topology

import (
	"fmt"

	"github.com/yannickkirschen/graphs"
	"github.com/yannickkirschen/graphs/inventory"
)

type Topology[O, C, P comparable] struct {
	inv   *inventory.Inventory[O, C, P]
	graph *graphs.Graph[O, P]
}

func NewTopology[O, C, P comparable]() *Topology[O, C, P] {
	return &Topology[O, C, P]{nil, nil}
}

func (top *Topology[O, C, P]) FindRef(fromRef, toRef O) ([][]*graphs.PathSegment[O, P], error) {
	from, err := top.inv.GetObject(fromRef).Take()
	if err != nil {
		return nil, fmt.Errorf("from ref %v not found in inventory: %s", fromRef, err)
	}

	to, err := top.inv.GetObject(toRef).Take()
	if err != nil {
		return nil, fmt.Errorf("to ref %v not found in inventory: %s", toRef, err)
	}

	if from.Class.PathConstruction == nil || from.Class.PathConstruction.Start == nil {
		return nil, fmt.Errorf("object %s (ID %v) does not allow paths to start here", from.Label, from.Id)
	}

	if to.Class.PathConstruction == nil || to.Class.PathConstruction.End == nil {
		return nil, fmt.Errorf("object %s (ID %v) does not allow paths to end here", to.Label, to.Id)
	}

	fromPort := from.Class.PathConstruction.Start
	toPort := to.Class.PathConstruction.End

	return top.graph.FindRef(from.Id, fromPort.Id, to.Id, toPort.Id)
}

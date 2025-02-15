package topology

import (
	"fmt"

	"github.com/yannickkirschen/graphs"
	"github.com/yannickkirschen/graphs/inventory"
)

type Topology struct {
	inv   *inventory.Inventory
	graph *graphs.Graph[inventory.Id, inventory.Id]
}

func NewTopology() *Topology {
	return &Topology{nil, nil}
}

func (top *Topology) FindRef(fromRef, toRef inventory.Id) ([][]*graphs.Triple[*inventory.Id, *graphs.Node[inventory.Id, inventory.Id], *inventory.Id], error) {
	from, err := top.inv.GetObject(fromRef).Take()
	if err != nil {
		return nil, fmt.Errorf("from ref %d not found in inventory: %s", fromRef, err)
	}

	to, err := top.inv.GetObject(toRef).Take()
	if err != nil {
		return nil, fmt.Errorf("to ref %d not found in inventory: %s", toRef, err)
	}

	if from.Class.PathConstruction == nil || from.Class.PathConstruction.Start == nil {
		return nil, fmt.Errorf("object %s (ID %d) does not allow paths to start here", from.Label, from.Id)
	}

	if to.Class.PathConstruction == nil || to.Class.PathConstruction.End == nil {
		return nil, fmt.Errorf("object %s (ID %d) does not allow paths to end here", to.Label, to.Id)
	}

	fromPort := from.Class.PathConstruction.Start
	toPort := to.Class.PathConstruction.End

	return top.graph.FindRef(from.Id, fromPort.Id, to.Id, toPort.Id)
}

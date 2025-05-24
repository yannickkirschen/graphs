package topology

import (
	"fmt"
	"io"
	"os"

	"github.com/yannickkirschen/graphs"
	"github.com/yannickkirschen/graphs/inventory"
	"gopkg.in/yaml.v3"
)

type Model struct {
	Connections []*Connection `yaml:"connections"`
}

type Connection struct {
	From          inventory.Id `yaml:"from"`
	FromPort      inventory.Id `yaml:"fromPort"`
	To            inventory.Id `yaml:"to"`
	ToPort        inventory.Id `yaml:"toPort"`
	Bidirectional bool         `yaml:"bidirectional"`
}

func (model *Model) ToGraph(inv *inventory.Inventory) *graphs.Graph[inventory.Id, inventory.Id] {
	g := graphs.NewGraph[inventory.Id, inventory.Id]()
	for _, object := range inv.Objects() {
		g.AddNode(object.ToGraphNode())
	}

	for _, connection := range model.Connections {
		if connection.Bidirectional {
			g.ConnectRefBi(
				inventory.Id(connection.From),
				inventory.Id(connection.FromPort),
				inventory.Id(connection.To),
				inventory.Id(connection.ToPort),
			)
		} else {
			g.ConnectRef(
				inventory.Id(connection.From),
				inventory.Id(connection.FromPort),
				inventory.Id(connection.To),
				inventory.Id(connection.ToPort),
			)
		}
	}

	return g
}

func (model *Model) ToTopology(inv *inventory.Inventory) *Topology {
	return &Topology{inv, model.ToGraph(inv)}
}

func Parse(inv *inventory.Inventory, r io.ReadCloser) (*Topology, error) {
	var model *Model
	if err := yaml.NewDecoder(r).Decode(&model); err != nil {
		return nil, fmt.Errorf("error parsing input: %s", err)
	}

	return model.ToTopology(inv), nil
}

func ParseFile(inv *inventory.Inventory, filename string) (*Topology, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error opening %s: %s", filename, err)
	}

	return Parse(inv, f)
}

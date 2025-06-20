package topology

import (
	"fmt"
	"io"
	"os"

	"github.com/yannickkirschen/graphs"
	"github.com/yannickkirschen/graphs/inventory"
	"gopkg.in/yaml.v3"
)

type Model[O, C, P comparable] struct {
	Connections []*Connection[O, C, P] `yaml:"connections"`
}

type Connection[O, C, P comparable] struct {
	From          O    `yaml:"from"`
	FromPort      P    `yaml:"fromPort"`
	To            O    `yaml:"to"`
	ToPort        P    `yaml:"toPort"`
	Bidirectional bool `yaml:"bidirectional"`
}

func (model *Model[O, C, P]) ToGraph(inv *inventory.Inventory[O, C, P]) *graphs.Graph[O, P] {
	g := graphs.NewGraph[O, P]()
	for _, object := range inv.Objects() {
		g.AddNode(object.ToGraphNode())
	}

	for _, connection := range model.Connections {
		if connection.Bidirectional {
			g.ConnectRefBi(
				connection.From,
				connection.FromPort,
				connection.To,
				connection.ToPort,
			)
		} else {
			g.ConnectRef(
				connection.From,
				connection.FromPort,
				connection.To,
				connection.ToPort,
			)
		}
	}

	return g
}

func (model *Model[O, C, P]) ToTopology(inv *inventory.Inventory[O, C, P]) *Topology[O, C, P] {
	return &Topology[O, C, P]{inv, model.ToGraph(inv)}
}

func Parse[O, C, P comparable](inv *inventory.Inventory[O, C, P], r io.ReadCloser) (*Topology[O, C, P], error) {
	var model *Model[O, C, P]
	if err := yaml.NewDecoder(r).Decode(&model); err != nil {
		return nil, fmt.Errorf("error parsing input: %s", err)
	}

	return model.ToTopology(inv), nil
}

func ParseFile[O, C, P comparable](inv *inventory.Inventory[O, C, P], filename string) (*Topology[O, C, P], error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error opening %s: %s", filename, err)
	}

	return Parse(inv, f)
}

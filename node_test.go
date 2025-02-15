package graphs_test

import (
	"testing"

	"github.com/yannickkirschen/graphs"
)

// The tested node is a railway point and looks like this:
//
// head ---------- main
//         `------ diversion
//
// We're testing a point called "W1".

func MakeNode() *graphs.Node[string, string] {
	node := graphs.NewNode[string, string]("W1")

	node.ConnectBi("head", "main")
	node.ConnectBi("head", "diversion")

	return node
}

func TestId(t *testing.T) {
	node := MakeNode()
	if node.Id() != "W1" {
		t.Fatalf("expected node ID to be 'W1', but got %v", node.Id())
	}
}

func TestHead(t *testing.T) {
	node := MakeNode()

	nextPorts := node.Next("head")
	if len(nextPorts) != 2 {
		t.Fatalf("expected 2 ports, but got %d: %v", len(nextPorts), nextPorts)
	}

	if !(nextPorts[0] == "main" && nextPorts[1] == "diversion") && !(nextPorts[0] == "diversion" && nextPorts[1] == "main") {
		t.Fatalf("expected next ports to be ['main', 'diversion'], but got %v", nextPorts)
	}
}

func TestMain(t *testing.T) {
	node := MakeNode()

	nextPorts := node.Next("main")
	if len(nextPorts) != 1 {
		t.Fatalf("expected 1 ports, but got %d: %v", len(nextPorts), nextPorts)
	}

	if nextPorts[0] != "head" {
		t.Fatalf("expected next ports to be ['head'], but got %v", nextPorts)
	}
}

func TestDiversion(t *testing.T) {
	node := MakeNode()

	nextPorts := node.Next("diversion")
	if len(nextPorts) != 1 {
		t.Fatalf("expected 1 ports, but got %d: %v", len(nextPorts), nextPorts)
	}

	if nextPorts[0] != "head" {
		t.Fatalf("expected next ports to be ['head'], but got %v", nextPorts)
	}
}

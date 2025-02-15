package graphs_test

import (
	"testing"

	"github.com/yannickkirschen/graphs"
)

/*
       e a,-4-,bc
    1 - 2 - 3 - 5 - 6
   a b c d a b b a a b
*/

func TestFind(t *testing.T) {
	one := graphs.NewNode[int, string](1)
	two := graphs.NewNode[int, string](2)
	three := graphs.NewNode[int, string](3)
	four := graphs.NewNode[int, string](4)
	five := graphs.NewNode[int, string](5)
	six := graphs.NewNode[int, string](6)

	one.ConnectBi("a", "b")
	two.ConnectBi("c", "d")
	two.ConnectBi("c", "e")
	three.ConnectBi("a", "b")
	four.ConnectBi("a", "b")
	five.ConnectBi("a", "b")
	five.ConnectBi("a", "c")
	six.ConnectBi("a", "b")

	graph := graphs.NewGraph[int, string]()
	graph.AddNode(one)
	graph.AddNode(two)
	graph.AddNode(three)
	graph.AddNode(four)
	graph.AddNode(five)
	graph.AddNode(six)

	graph.ConnectRefBi(1, "b", 2, "c")
	graph.ConnectRefBi(2, "d", 3, "a")
	graph.ConnectRefBi(2, "e", 4, "a")
	graph.ConnectRefBi(3, "b", 5, "b")
	graph.ConnectRefBi(4, "b", 5, "c")
	graph.ConnectRefBi(5, "a", 6, "a")

	paths, err := graph.FindRef(1, "b", 4, "b")
	if err != nil {
		t.Fatalf("error when finding paths: %s", err)
	}

	if len(paths) != 1 {
		t.Fatalf("expected 1 path, but got %d: %v", len(paths), paths)
	}

	path := paths[0]
	if len(path) != 3 {
		t.Fatalf("path has wrong length: expected 3, but got %d: %v", len(path), path)
	}

	if path[0].Middle.Id() != 1 || path[1].Middle.Id() != 2 || path[2].Middle.Id() != 4 {
		t.Fatalf("expected path to be 1 -> 2 -> 4 but got %v", path)
	}
}

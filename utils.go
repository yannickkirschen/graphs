package graphs

import "fmt"

type Triple[L, M, R comparable] struct {
	Left   L
	Middle M
	Right  R
}

func NewTriple[L, M, R comparable](l L, m M, r R) *Triple[L, M, R] {
	return &Triple[L, M, R]{l, m, r}
}

func (triple *Triple[L, M, R]) String() string {
	return fmt.Sprintf("<%v, %v, %v>", triple.Left, triple.Middle, triple.Right)
}

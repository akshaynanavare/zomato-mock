package utils

import (
	hp "container/heap"
)

type Path struct {
	Value float64
	Nodes []string
}

type minPath []Path

func (h minPath) Len() int           { return len(h) }
func (h minPath) Less(i, j int) bool { return h[i].Value < h[j].Value }
func (h minPath) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *minPath) Push(x interface{}) {
	*h = append(*h, x.(Path))
}

func (h *minPath) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

type heap struct {
	values *minPath
}

func NewHeap() *heap {
	return &heap{values: &minPath{}}
}

func (h *heap) Push(p Path) {
	hp.Push(h.values, p)
}

func (h *heap) Pop() Path {
	i := hp.Pop(h.values)
	return i.(Path)
}

func (h *heap) Len() int {
	return len(*h.values)
}

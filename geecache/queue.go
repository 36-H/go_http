package geecache

import (
	"container/heap"
)

type lfu_entry struct {
	entry
	weight int
	index  int
}

type queue []*lfu_entry

func (q queue) Len() int {
	return len(q)
}

func (q queue) Less(i, j int) bool {
	return q[i].weight < q[j].weight
}

func (q queue) Swap(i, j int) {
	q[i], q[j] = q[j], q[i]
	q[i].index = i
	q[j].index = j
}

func (q *queue) Push(x interface{}) {
	n := len(*q)
	en := x.(*lfu_entry)
	en.index = n
	*q = append(*q, en)
}

func (q *queue) Pop() interface{} {
	old := *q
	n := len(old)
	en := old[n-1]
	old[n-1] = nil // avoid memory leak
	en.index = -1  // for safety
	*q = old[0 : n-1]
	return en
}

func (q *queue) update(en *lfu_entry, value interface{}, weight int) {
	en.Value = value
	en.weight = weight
	heap.Fix(q, en.index)
}
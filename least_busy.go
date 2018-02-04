package slb

import (
	"bytes"
	"container/heap"
	"fmt"
	"strings"
)

type leastBusy nodes

func newLeastBusy(urls []string) leastBusy {
	pool := make(leastBusy, len(urls))
	for i, url := range urls {
		pool[i] = newNode(url, i)
	}
	heap.Init(&pool)
	return pool
}

func (p leastBusy) Dispatch() *node {
	node := heap.Pop(&p).(*node)
	node.pending += 1
	heap.Push(&p, node)
	heap.Fix(&p, node.index)
	return node
}

func (p leastBusy) Complete(host string) {
	var n *node
	for _, node := range p {
		if strings.Compare(node.host, host) == 0 {
			n = node
			break
		}
	}

	n.pending -= 1
	heap.Fix(&p, n.index)
}

func (p leastBusy) Len() int {
	return len(p)
}

func (p leastBusy) Less(i, j int) bool {
	return p[i].pending < p[j].pending
}

func (p leastBusy) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
	p[i].index = i
	p[j].index = j
}

func (p *leastBusy) Push(x interface{}) {
	n := x.(*node)
	n.index = p.Len()
	*p = append(*p, n)
}

func (p *leastBusy) Pop() interface{} {
	pool := *p
	last := p.Len() - 1
	node := pool[last]
	*p = pool[:last]
	return node
}

func (p leastBusy) String() string {
	var output bytes.Buffer
	output.WriteString("\nHost with pending tasks: \n")
	for _, node := range p {
		str := fmt.Sprintf("%+v\n", *node)
		output.WriteString(str)
	}

	output.WriteString("\n")

	mean := mean(nodes(p))
	stdDev := standardDeviation(nodes(p))

	stats := fmt.Sprintf("Avg Load: %.2f | Std Dev: %.2f\n", mean, stdDev)
	output.WriteString(stats)

	return output.String()
}

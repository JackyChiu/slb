package slb

import (
	"bytes"
	"fmt"
	"math"
)

type node struct {
	host    string
	index   int
	pending int
}

func newNode(url string, index int) *node {
	return &node{
		host:    url,
		pending: 0,
		index:   index,
	}
}

type nodes []*node

func newNodes(urls []string) nodes {
	pool := make(nodes, len(urls))
	for i, url := range urls {
		pool[i] = newNode(url, i)
	}
	return pool
}

func (n nodes) Len() int {
	return len(n)
}

func (n nodes) Less(i, j int) bool {
	return n[i].pending < n[j].pending
}

func (n nodes) Swap(i, j int) {
	n[i], n[j] = n[j], n[i]
	n[i].index = i
	n[j].index = j
}

func (n *nodes) Push(x interface{}) {
	node := x.(*node)
	node.index = n.Len()
	*n = append(*n, node)
}

func (n *nodes) Pop() interface{} {
	nodes := *n
	last := n.Len() - 1
	node := nodes[last]
	*n = nodes[:last]
	return node
}

func (n nodes) String() string {
	var output bytes.Buffer
	output.WriteString("\nHost with pending tasks: \n")
	for _, node := range n {
		str := fmt.Sprintf("%+v\n", *node)
		output.WriteString(str)
	}

	output.WriteString("\n")

	mean := n.mean()
	stdDev := n.standardDeviation()

	stats := fmt.Sprintf("Avg Load: %.2f | Std Dev: %.2f\n", mean, stdDev)
	output.WriteString(stats)

	return output.String()
}

func (n nodes) mean() (mean float64) {
	length := float64(len(n))
	for _, server := range n {
		mean += float64(server.pending)
	}
	return mean / length
}

func (n nodes) standardDeviation() (stdDev float64) {
	length := float64(len(n))
	mean := n.mean()
	for _, server := range n {
		stdDev += math.Pow((float64(server.pending) - mean), 2)
	}
	return math.Sqrt((1 / length) * stdDev)
}

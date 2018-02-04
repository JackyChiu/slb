package slb

import (
	"container/heap"
	"strings"
)

type leastBusy struct {
	nodes nodes
}

func newLeastBusy(urls []string) leastBusy {
	nodes := newNodes(urls)
	heap.Init(&nodes)
	return leastBusy{
		nodes: nodes,
	}
}

func (p leastBusy) Dispatch() *node {
	node := heap.Pop(&p.nodes).(*node)
	node.pending += 1
	heap.Push(&p.nodes, node)
	heap.Fix(&p.nodes, node.index)
	return node
}

func (p leastBusy) Complete(host string) {
	var n *node
	for _, node := range p.nodes {
		if strings.Compare(node.host, host) == 0 {
			n = node
			break
		}
	}

	n.pending -= 1
	heap.Fix(&p.nodes, n.index)
}

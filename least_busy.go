package slb

import (
	"container/heap"
	"net/http"
	"strings"
)

// leastBusy is a scheduling strategy that looks at the amount of work
// that the workers have and schedules jobs to the least busy worker.
type leastBusy struct {
	nodes        nodes
	dispatchChan chan chan node
	completeChan chan *http.Response
}

func newLeastBusy(hosts []string) *leastBusy {
	nodes := newNodes(hosts)
	heap.Init(&nodes)
	lb := &leastBusy{
		nodes:        nodes,
		dispatchChan: make(chan chan node),
		completeChan: make(chan *http.Response),
	}
	go lb.balance()
	return lb
}

func (l *leastBusy) Dispatch() node {
	nodeChan := make(chan node)
	l.dispatchChan <- nodeChan
	return <-nodeChan
}

func (l *leastBusy) Complete(res *http.Response) {
	l.completeChan <- res
}

func (l *leastBusy) balance() {
	for {
		select {
		case nodeChan := <-l.dispatchChan:
			nodeChan <- l.dispatch()
		case res := <-l.completeChan:
			l.complete(res.Request.URL.Host)
		}
	}
}

func (l *leastBusy) dispatch() node {
	node := heap.Pop(&l.nodes).(*node)
	node.pending += 1
	heap.Push(&l.nodes, node)
	heap.Fix(&l.nodes, node.index)
	return *node
}

func (l *leastBusy) complete(host string) {
	var n *node
	for _, node := range l.nodes {
		if strings.Compare(node.host, host) == 0 {
			n = node
			break
		}
	}

	n.pending -= 1
	heap.Fix(&l.nodes, n.index)
}

func (l *leastBusy) String() string {
	return l.nodes.String()
}

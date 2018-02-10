package slb

import (
	"container/ring"
	"net/http"
	"strings"
)

type roundRobin struct {
	ring         *ring.Ring
	nodes        nodes
	dispatchChan chan chan node
	completeChan chan *http.Response
}

func newRoundRobin(hosts []string) *roundRobin {
	nodes := newNodes(hosts)
	r := ring.New(len(hosts))
	for i, node := range nodes {
		node.index = i
		r.Value = node
		r = r.Next()
	}
	rr := &roundRobin{
		ring:         r,
		nodes:        nodes,
		dispatchChan: make(chan chan node),
		completeChan: make(chan *http.Response),
	}
	go rr.balance()
	return rr
}

func (rr *roundRobin) Dispatch() <-chan node {
	nodeChan := make(chan node)
	rr.dispatchChan <- nodeChan
	return nodeChan
}

func (rr *roundRobin) Complete(res *http.Response) {
	rr.completeChan <- res
}

func (rr *roundRobin) balance() {
	for {
		select {
		case nodeChan := <-rr.dispatchChan:
			nodeChan <- rr.dispatch()
		case res := <-rr.completeChan:
			rr.complete(res.Request.URL.Host)
		}
	}
}

func (rr *roundRobin) dispatch() node {
	node := rr.ring.Value.(*node)
	node.pending += 1
	rr.ring = rr.ring.Next()
	return *node
}

func (rr *roundRobin) complete(host string) {
	rr.ring.Do(func(value interface{}) {
		node := value.(*node)
		if strings.Compare(node.host, host) == 0 {
			node.pending -= 1
		}
	})
}

func (rr *roundRobin) String() string {
	return rr.nodes.String()
}

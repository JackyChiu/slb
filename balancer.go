package slb

import (
	"log"
	"net/http"
	"net/http/httputil"

	"github.com/pkg/errors"
)

// Available strategies for the pool.
const (
	LeastBusy  = "least_busy"
	RoundRobin = "round_robin"
)

// Pool is an interface for pools with different
// strategies of distributing work.
type Pool interface {
	Dispatch() <-chan node
	Complete(res *http.Response)
}

// NewPool provides a pool with specified strategy.
func NewPool(strategy string, hosts []string) Pool {
	switch strategy {
	case LeastBusy:
		return newLeastBusy(hosts)
	case RoundRobin:
		return newRoundRobin(hosts)
	default:
		panic(errors.Errorf("%v is not a valid strategy", strategy))
	}
}

// Balancer is the reverse proxy server that balances requests.
type Balancer struct {
	*httputil.ReverseProxy
	pool Pool
}

// NewBalancer creates a new balancer to balance requests between hosts
// and uses specified strategy.
func NewBalancer(strategy string, hosts []string) *Balancer {
	b := &Balancer{
		pool: NewPool(strategy, hosts),
	}

	b.ReverseProxy = &httputil.ReverseProxy{
		Director:       b.Director,
		ModifyResponse: b.ModifyResponse,
	}

	return b
}

// Director directs the request to the node that was dispatched by pool.
func (b *Balancer) Director(r *http.Request) {
	nodeChan := b.pool.Dispatch()
	node := <-nodeChan
	log.Println(b.pool)

	r.URL.Scheme = "http"
	r.URL.Host = node.host
}

// ModifyResponse tells the pool that the request was handled.
func (b *Balancer) ModifyResponse(res *http.Response) error {
	b.pool.Complete(res)
	return nil
}

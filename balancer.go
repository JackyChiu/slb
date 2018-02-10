package slb

import (
	"log"
	"net/http"
	"net/http/httputil"

	"github.com/pkg/errors"
)

const (
	LeastBusy  = "least_busy"
	RoundRobin = "round_robin"
)

type Pool interface {
	Dispatch() <-chan node
	Complete(res *http.Response)
}

func NewPool(strategy string, hosts []string) Pool {
	switch strategy {
	case LeastBusy:
		return newLeastBusy(hosts)
	case RoundRobin:
		return newRoundRobin(hosts)
	default:
		panic(errors.Errorf("%v is not a valid stratgey", strategy))
	}
}

type Balancer struct {
	*httputil.ReverseProxy
	pool Pool
}

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

func (b *Balancer) Director(r *http.Request) {
	nodeChan := b.pool.Dispatch()
	node := <-nodeChan
	log.Println(b.pool)

	r.URL.Scheme = "http"
	r.URL.Host = node.host
}

func (b *Balancer) ModifyResponse(res *http.Response) error {
	b.pool.Complete(res)
	return nil
}

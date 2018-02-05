package slb

import (
	"log"
	"net/http"
	"net/http/httputil"
)

type Pool interface {
	Dispatch() <-chan node
	Complete(res *http.Response)
}

func NewPool(hosts []string) Pool {
	return newLeastBusy(hosts)
}

type Balancer struct {
	*httputil.ReverseProxy
	pool Pool
}

func NewBalancer(hosts []string) *Balancer {
	b := &Balancer{
		pool: NewPool(hosts),
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

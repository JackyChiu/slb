package slb

import (
	"log"
	"net/http"
	"net/http/httputil"
)

type Pool interface {
	Node(host string) *node
	Dispatch() *node
	Complete()
}

type Balancer struct {
	*httputil.ReverseProxy
	pool leastBusy
}

func NewBalancer(hosts []string) *Balancer {
	b := &Balancer{
		pool: newLeastBusy(hosts),
	}

	b.ReverseProxy = &httputil.ReverseProxy{
		Director:       b.Director,
		ModifyResponse: b.ModifyResponse,
	}

	return b
}

func (b *Balancer) Director(r *http.Request) {
	node := b.pool.Dispatch()
	log.Println(b.pool)

	r.URL.Scheme = "http"
	r.URL.Host = node.host
}

func (b *Balancer) ModifyResponse(res *http.Response) error {
	b.pool.Complete(res.Request.URL.Host)
	return nil
}

package slb

import (
	"container/heap"
	"log"
	"net/http"
	"net/http/httputil"

	"github.com/pkg/errors"
)

type Balancer struct {
	*httputil.ReverseProxy
	pool pool
}

func NewBalancer(hosts []string) *Balancer {
	b := &Balancer{
		pool: newPool(hosts),
	}

	b.ReverseProxy = &httputil.ReverseProxy{
		Director:       b.Director,
		ModifyResponse: b.ModifyResponse,
	}

	return b
}

func (b *Balancer) Director(r *http.Request) {
	server := b.pool.Pop().(*server)
	log.Printf("serving to %v", server.host)
	server.pending += 1

	r.URL.Scheme = "http"
	r.URL.Host = server.host

	b.pool.Push(server)
	heap.Fix(&b.pool, server.index)
}

func (b *Balancer) ModifyResponse(res *http.Response) error {
	server, err := b.pool.Server(res.Request.URL.Host)
	if err != nil {
		return errors.Wrap(err, "couldn't update pool")
	}
	server.pending -= 1
	heap.Fix(&b.pool, server.index)
	return nil
}

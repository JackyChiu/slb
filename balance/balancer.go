package balance

import (
	"container/heap"
	"net/http"
	"net/http/httputil"

	"github.com/pkg/errors"
)

type Balancer struct {
	*httputil.ReverseProxy
	pool pool
}

func NewBalancer() *Balancer {
	b := &Balancer{
		pool: newPool([]string{
			"localhost:9000",
			"localhost:9001",
			"localhost:9002",
			"localhost:9003",
			"localhost:9004",
		}),
	}

	b.ReverseProxy = &httputil.ReverseProxy{
		Director:       b.Director,
		ModifyResponse: b.ModifyResponse,
	}

	return b
}

func (b *Balancer) Director(r *http.Request) {
	server := b.pool.Pop().(*server)
	server.pending += 1

	r.URL.Scheme = "http"
	r.URL.Host = server.host

	b.pool.Push(server)
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

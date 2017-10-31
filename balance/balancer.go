package balance

import (
	"log"
	"net/http"
	"net/http/httputil"
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
	//worker := b.pool.Pop().(*worker)
	//worker.pending += 1
	worker := b.pool[0]
	r.URL.Scheme = "http"
	r.URL.Host = worker.host
}

func (b *Balancer) ModifyResponse(res *http.Response) error {
	log.Printf("response %+v", res)
	// TODO I know what server from this?
	log.Printf("host from res %s", res.Request.URL.Host)
	return nil
}

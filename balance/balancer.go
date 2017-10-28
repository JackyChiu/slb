package balance

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type Balancer struct {
	*httputil.ReverseProxy
	pool pool
}

func NewBalancer() *Balancer {
	b := &Balancer{
		pool: newPool([]string{
			"http://localhost:9000",
			"http://localhost:9001",
			"http://localhost:9002",
			"http://localhost:9003",
			"http://localhost:9004",
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
	// TODO: halp, how do I do this cleanly
	// Read: https://golang.org/src/net/http/httputil/reverseproxy.go#L61
	u, _ := url.Parse(worker.host)
	r.URL = u
}

func (b *Balancer) ModifyResponse(res *http.Response) error {
	log.Printf("%+v", res)
	return nil
}

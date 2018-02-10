package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/JackyChiu/slb"
)

func main() {
	var (
		configPath = flag.String("config", "", "path to config file")
		strategy   = flag.String("strategy", slb.RoundRobin, "strategy of the load balancer")
	)
	flag.Parse()

	config := slb.MustParseConfig(*configPath)
	port := fmt.Sprintf(":%v", config.Port)

	log.Printf("balancing from port %v", port)
	http.ListenAndServe(port, slb.NewBalancer(*strategy, config.Hosts))
}

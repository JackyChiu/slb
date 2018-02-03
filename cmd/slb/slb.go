package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/JackyChiu/slb"
)

func main() {
	var (
		configPath = flag.String("c", "", "path to config file")
	)
	flag.Parse()

	config := slb.MustParseConfig(*configPath)
	port := fmt.Sprintf(":%v", config.Port)

	http.ListenAndServe(port, slb.NewBalancer(config.Hosts))
}

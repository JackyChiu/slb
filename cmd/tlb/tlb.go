package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/JackyChiu/tlb"
	"github.com/JackyChiu/tlb/balance"
)

func main() {
	var (
		configPath = flag.String("c", "", "path to config file")
	)
	flag.Parse()

	config := tlb.MustParseConfig(*configPath)
	port := fmt.Sprintf(":%v", config.Port)

	http.ListenAndServe(port, balance.NewBalancer(config.Hosts))
}

package main

import (
	"net/http"

	"github.com/JackyChiu/tlb/balance"
	"github.com/JackyChiu/tlb/server"
)

func main() {
	// Run the servers
	server.StartServers(server.DEFAULT_PORTS)

	// Run frontend server
	http.ListenAndServe(":8000", balance.NewBalancer())

	//var buf bytes.Buffer
	//buf.WriteString(`{ "seconds": 3 }`)
	//res, err := http.Post("http://localhost:8000", "application/json", &buf)
	//if err != nil {
	//	log.Fatalf("req failed: %v", err)
	//}
	//log.Printf("yay it works: %v", res)
}

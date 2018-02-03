package main

import (
	"log"
	"net/http"
	"strings"
)

func main() {
	body := strings.NewReader(`{ "seconds": 3 }`)
	res, err := http.Post("http://localhost:8000", "application/json", body)
	if err != nil {
		log.Fatalf("req failed: %v", err)
	}
	log.Printf("yay it works: %v", res)
}

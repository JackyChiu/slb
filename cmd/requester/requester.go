package main

import (
	"bytes"
	"log"
	"net/http"
)

func main() {
	var buf bytes.Buffer
	buf.WriteString(`{ "seconds": 3 }`)
	res, err := http.Post("http://localhost:8000", "application/json", &buf)
	if err != nil {
		log.Fatalf("req failed: %v", err)
	}
	log.Printf("yay it works: %v", res)
}

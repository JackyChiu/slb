package main

import (
	"encoding/json"
	"errors"
	"flag"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/JackyChiu/slb"
)

func main() {
	var (
		configPath = flag.String("c", "", "path to config file")
	)
	flag.Parse()

	config := slb.MustParseConfig(*configPath)
	var ports []string
	for _, host := range config.Hosts {
		_, port, err := net.SplitHostPort(host)
		if err != nil {
			panic(errors.New("server only runs on localhost through ports"))
		}
		ports = append(ports, ":"+port)
	}

	http.HandleFunc("/", sleepHandler)

	var wg sync.WaitGroup
	for _, port := range ports {
		port := port
		wg.Add(1)

		go func() {
			defer wg.Done()
			log.Printf("Starting up server on %v", port)
			StartServer(port)
		}()
	}

	wg.Wait()
}

// StartServer starts up a indiviual server running on specified port
func StartServer(port string) {
	http.ListenAndServe(port, nil)
}

var (
	defaultSleep = 25 * time.Second
)

type sleepRequest struct {
	Seconds int `json:"seconds"`
}

// sleepHandler simulates a route that does work by sleeping
func sleepHandler(w http.ResponseWriter, r *http.Request) {
	var req sleepRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "unexpected request format", http.StatusUnprocessableEntity)
		return
	}
	defer r.Body.Close()

	log.Printf("%v recieved a request to sleep for %v seconds", r.URL.Host, req.Seconds)
	select {
	case <-time.After(time.Duration(req.Seconds) * time.Second):
	case <-time.After(defaultSleep):
	}

	w.WriteHeader(http.StatusOK)
}

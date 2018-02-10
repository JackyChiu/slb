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
		configPath = flag.String("config", "", "path to config file")
	)
	flag.Parse()

	config := slb.MustParseConfig(*configPath)

	var wg sync.WaitGroup
	for _, host := range config.Hosts {
		host := host
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, port, err := net.SplitHostPort(host)
			if err != nil {
				panic(errors.New("server only runs on localhost through ports"))
			}

			log.Printf("Starting up server on %v", port)
			StartServer(host, ":"+port)
		}()
	}

	wg.Wait()
}

// StartServer starts up a indiviual server running on specified port.
func StartServer(host, port string) {
	server := http.Server{
		Addr: port,
		Handler: &sleepHandler{
			host: host,
		},
	}
	log.Fatal(server.ListenAndServe())
}

var maxSleep = 25 * time.Second

type (
	// sleepHandler simulates a route that does work by sleeping.
	sleepHandler struct {
		host string
	}

	// sleepRequest defines for long to sleep
	sleepRequest struct {
		Seconds int `json:"seconds"`
	}
)

func (s *sleepHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req sleepRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "unexpected request format", http.StatusUnprocessableEntity)
		return
	}
	defer r.Body.Close()

	log.Printf("%v recieved a request to sleep for %v seconds", s.host, req.Seconds)
	select {
	case <-time.After(time.Duration(req.Seconds) * time.Second):
	case <-time.After(maxSleep):
	}

	w.WriteHeader(http.StatusOK)
}

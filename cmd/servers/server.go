package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"
)

func main() {
	http.HandleFunc("/", sleepHandler)
	ports := []string{
		":9000",
		":9001",
		":9002",
		":9003",
	}

	var wg sync.WaitGroup
	for _, port := range ports {
		port := port
		wg.Add(1)

		go func() {
			defer wg.Done()
			StartServer(port)
		}()
	}

	wg.Wait()
}

// StartServer starts up a indiviual server running on specified port
func StartServer(port string) {
	http.ListenAndServe(port, nil)
}

// StartServers starts up a bunch of servers listening given ports
func StartServers(ports []string) {
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

	log.Printf("%v recieved a request to sleep for %v seconds", r.Host, req.Seconds)
	select {
	case <-time.After(time.Duration(req.Seconds) * time.Second):
	case <-time.After(defaultSleep):
	}

	w.WriteHeader(http.StatusOK)
}

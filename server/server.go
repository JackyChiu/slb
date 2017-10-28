package server

import (
	"encoding/json"
	"net/http"
	"time"
)

var (
	// DEFAULT_PORTS are the default ports to listen too when starting new servers
	DEFAULT_PORTS = []string{
		":9000",
		":9001",
		":9002",
		":9003",
		":9004",
	}
	defaultSleep = 25 * time.Second
)

// StartServer starts up a indiviual server running on specified port
func StartServer(port string) {
	http.ListenAndServe(port, nil)
}

// StartServers starts up a bunch of servers listening given ports
func StartServers(ports []string) {
	http.HandleFunc("/sleep", sleepHandler)

	for _, port := range ports {
		go StartServer(port)
	}
}

type sleepRequest struct {
	duration time.Duration `json:"duration"`
}

// sleepHandler simulates a route that does work by sleeping
func sleepHandler(w http.ResponseWriter, r *http.Request) {
	var req *sleepRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "unexpected request format", http.StatusUnprocessableEntity)
	}

	select {
	case <-time.After(req.duration):
	case <-time.After(defaultSleep):
	}

	w.WriteHeader(http.StatusOK)
}

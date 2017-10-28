package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
)

func TestStartServers_Default(t *testing.T) {
	var sleepReq = sleepRequest{
		Seconds: 1,
	}

	StartServers(DEFAULT_PORTS)

	for _, port := range DEFAULT_PORTS {
		t.Run("testing "+port, func(t *testing.T) {
			t.Parallel()
			var buf bytes.Buffer
			json.NewEncoder(&buf).Encode(sleepReq)
			res, err := http.Post("http://localhost"+port, "application/json", &buf)
			if err != nil {
				t.Fatalf("unexpected response: %s, err: %s", res, err)
			}
		})
	}
}

package main

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

	StartServers(DefaultPorts)

	for _, port := range DefaultPorts {
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

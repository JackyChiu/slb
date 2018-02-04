package slb

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBalancer(t *testing.T) {
	expected := "Hello, client"
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, expected)
	}))
	defer testServer.Close()

	balancer := NewBalancer([]string{testServer.URL[7:]}) // hack to get host and port
	balancerServer := httptest.NewServer(balancer)

	res, err := http.Post(balancerServer.URL, "application/json", nil)
	require.NoError(t, err)

	actual, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	require.NoError(t, err)

	require.Equal(t, expected, string(actual))
}

package balance

import (
	"container/heap"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPool(t *testing.T) {
	var (
		testServers = []*server{
			&server{
				host:    "localhost:9000",
				pending: 7,
			},
			&server{
				host:    "localhost:9001",
				pending: 9,
			},
			&server{
				host:    "localhost:9002",
				pending: 3,
			},
			&server{
				host:    "localhost:9003",
				pending: 4,
			},
		}
		testPool = make(pool, len(testServers))
	)
	for i, server := range testServers {
		server.index = i
		testPool[i] = server
	}
	heap.Init(&testPool)

	for _, server := range testPool {
		fmt.Printf("%+v\n", server)
	}

	t.Run("Heap ordering", func(t *testing.T) {
		for i, expected := range []server{
			server{
				host:    "localhost:9002",
				pending: 3,
				index:   0,
			},
			server{
				host:    "localhost:9003",
				pending: 4,
				index:   1,
			},
			server{
				host:    "localhost:9001",
				pending: 7,
				index:   2,
			},
			server{
				host:    "localhost:9000",
				pending: 9,
				index:   3,
			},
		} {
			actual := *testPool[i]
			require.Equal(t, expected, actual, "unexpected value")
		}
	})
}

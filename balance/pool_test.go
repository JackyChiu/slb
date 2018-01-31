package balance

import (
	"container/heap"
	"testing"

	"github.com/stretchr/testify/require"
)

var hosts = []string{
	"localhost:9000",
	"localhost:9001",
	"localhost:9002",
	"localhost:9003",
	"localhost:9004",
}

func TestPool_newPool(t *testing.T) {
	testPool := newPool(hosts)
	expected := pool{
		{host: "localhost:9000", index: 0},
		{host: "localhost:9001", index: 1},
		{host: "localhost:9002", index: 2},
		{host: "localhost:9003", index: 3},
		{host: "localhost:9004", index: 4},
	}
	require.Equal(t, expected, testPool)
}

func TestPool_Len(t *testing.T) {
	testPool := newPool(hosts)
	require.Equal(t, len(hosts), testPool.Len())
}

func TestPool_Swap(t *testing.T) {
	testPool := newPool(hosts)
	testPool.Swap(1, 2)
	expected := pool{
		{host: "localhost:9000", index: 0},
		{host: "localhost:9002", index: 1},
		{host: "localhost:9001", index: 2},
		{host: "localhost:9003", index: 3},
		{host: "localhost:9004", index: 4},
	}
	require.Equal(t, expected, testPool)
}

func TestPool_Pop(t *testing.T) {
	testPool := newPool(hosts)

	actual := testPool.Pop()
	expected := &server{
		host:    "localhost:9004",
		pending: 0,
		index:   4,
	}
	require.Equal(t, expected, actual)

	testPool[3].pending = 2
	heap.Fix(&testPool, 3)

	actual = testPool.Pop()
	expected = &server{
		host:    "localhost:9001",
		pending: 0,
		index:   3,
	}
	require.Equal(t, expected, actual)
}

func TestPool_Push(t *testing.T) {
	testPool := newPool(hosts)

	testPool.Push(&server{
		host:    "new_host",
		pending: 2,
	})
	heap.Fix(&testPool, testPool.Len()-1)

	expected := pool{
		{host: "new_host", index: 0, pending: 2},
		{host: "localhost:9001", index: 1, pending: 0},
		{host: "localhost:9000", index: 2, pending: 0},
		{host: "localhost:9003", index: 3, pending: 0},
		{host: "localhost:9004", index: 4, pending: 0},
		{host: "localhost:9002", index: 5, pending: 0},
	}
	require.Equal(t, expected, testPool)
}

func TestPool_Server(t *testing.T) {
}

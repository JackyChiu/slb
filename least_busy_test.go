package slb

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

func TestLeastBusy_newLeastBusy(t *testing.T) {
	testPool := newLeastBusy(hosts)
	expected := leastBusy{
		{host: "localhost:9000", index: 0},
		{host: "localhost:9001", index: 1},
		{host: "localhost:9002", index: 2},
		{host: "localhost:9003", index: 3},
		{host: "localhost:9004", index: 4},
	}
	require.Equal(t, expected, testPool)
}

func TestLeastBusy_Len(t *testing.T) {
	testPool := newLeastBusy(hosts)
	require.Equal(t, len(hosts), testPool.Len())
}

func TestLeastBusy_Swap(t *testing.T) {
	testPool := newLeastBusy(hosts)
	testPool.Swap(1, 2)
	expected := leastBusy{
		{host: "localhost:9000", index: 0},
		{host: "localhost:9002", index: 1},
		{host: "localhost:9001", index: 2},
		{host: "localhost:9003", index: 3},
		{host: "localhost:9004", index: 4},
	}
	require.Equal(t, expected, testPool)
}

func TestLeastBusy_Pop(t *testing.T) {
	testPool := newLeastBusy(hosts)

	actual := testPool.Pop()
	expected := &node{
		host:    "localhost:9004",
		pending: 0,
		index:   4,
	}
	require.Equal(t, expected, actual)
}

func TestLeastBusy_Push(t *testing.T) {
	testPool := newLeastBusy(hosts)

	testPool.Push(&node{
		host:    "new_host",
		pending: 2,
	})

	expected := leastBusy{
		{host: "localhost:9000", index: 0, pending: 0},
		{host: "localhost:9001", index: 1, pending: 0},
		{host: "localhost:9002", index: 2, pending: 0},
		{host: "localhost:9003", index: 3, pending: 0},
		{host: "localhost:9004", index: 4, pending: 0},
		{host: "new_host", index: 5, pending: 2},
	}
	require.Equal(t, expected, testPool)
}

func TestLeastBusy_Dispatch(t *testing.T) {
	testPool := leastBusy{
		{host: "localhost:9000", index: 0, pending: 7},
		{host: "localhost:9002", index: 1, pending: 0},
		{host: "localhost:9001", index: 2, pending: 2},
		{host: "localhost:9003", index: 3, pending: 3},
		{host: "localhost:9004", index: 4, pending: 1},
	}
	heap.Init(&testPool)

	actual := testPool.Dispatch()
	expected := &node{host: "localhost:9002", index: 1, pending: 1}
	require.Equal(t, expected, actual)
}

func TestLeastBusy_Complete(t *testing.T) {
	testPool := leastBusy{
		{host: "localhost:9000", index: 0, pending: 7},
		{host: "localhost:9002", index: 1, pending: 0},
		{host: "localhost:9001", index: 2, pending: 2},
		{host: "localhost:9003", index: 3, pending: 3},
		{host: "localhost:9004", index: 4, pending: 1},
	}
	heap.Init(&testPool)

	testPool.Complete("localhost:9003")
	expected := leastBusy{
		{host: "localhost:9002", index: 0, pending: 0},
		{host: "localhost:9004", index: 1, pending: 1},
		{host: "localhost:9001", index: 2, pending: 2},
		{host: "localhost:9003", index: 3, pending: 2},
		{host: "localhost:9000", index: 4, pending: 7},
	}
	require.Equal(t, expected, testPool)
}

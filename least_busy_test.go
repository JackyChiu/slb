package slb

import (
	"container/heap"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLeastBusy_Dispatch(t *testing.T) {
	testPool := newLeastBusy([]string{})
	testPool.nodes = nodes{
		{host: "localhost:9000", index: 0, pending: 7},
		{host: "localhost:9002", index: 1, pending: 0},
		{host: "localhost:9001", index: 2, pending: 2},
		{host: "localhost:9003", index: 3, pending: 3},
		{host: "localhost:9004", index: 4, pending: 1},
	}
	heap.Init(&testPool.nodes)

	nodeChan := testPool.Dispatch()
	actual := <-nodeChan
	expected := node{host: "localhost:9002", index: 1, pending: 1}
	require.Equal(t, expected, actual)
}

func TestLeastBusy_complete(t *testing.T) {
	testPool := newLeastBusy([]string{})
	testPool.nodes = nodes{
		{host: "localhost:9000", index: 0, pending: 7},
		{host: "localhost:9002", index: 1, pending: 0},
		{host: "localhost:9001", index: 2, pending: 2},
		{host: "localhost:9003", index: 3, pending: 3},
		{host: "localhost:9004", index: 4, pending: 1},
	}
	heap.Init(&testPool.nodes)

	testPool.complete("localhost:9003")
	expected := nodes{
		{host: "localhost:9002", index: 0, pending: 0},
		{host: "localhost:9004", index: 1, pending: 1},
		{host: "localhost:9001", index: 2, pending: 2},
		{host: "localhost:9003", index: 3, pending: 2},
		{host: "localhost:9000", index: 4, pending: 7},
	}
	require.Equal(t, expected, testPool.nodes)
}

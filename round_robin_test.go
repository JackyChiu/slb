package slb

import (
	"container/ring"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRoundRobin_Dispatch(t *testing.T) {
	testPool := newRoundRobin([]string{})
	testNodes := nodes{
		{host: "localhost:9000", index: 0, pending: 7},
		{host: "localhost:9002", index: 1, pending: 0},
		{host: "localhost:9001", index: 2, pending: 2},
		{host: "localhost:9003", index: 3, pending: 3},
		{host: "localhost:9004", index: 4, pending: 1},
	}
	testPool.nodes = testNodes
	testPool.ring = ring.New(len(testNodes))

	for _, node := range testNodes {
		testPool.ring.Value = node
		testPool.ring = testPool.ring.Next()
	}

	actual := testPool.Dispatch()
	expected := node{host: "localhost:9000", index: 0, pending: 8}
	require.Equal(t, expected, actual)
}

func TestRoundRobin_complete(t *testing.T) {
	testPool := newRoundRobin([]string{})
	testNodes := nodes{
		{host: "localhost:9000", index: 0, pending: 7},
		{host: "localhost:9002", index: 1, pending: 0},
		{host: "localhost:9001", index: 2, pending: 2},
		{host: "localhost:9003", index: 3, pending: 3},
		{host: "localhost:9004", index: 4, pending: 1},
	}
	testPool.nodes = testNodes
	testPool.ring = ring.New(len(testNodes))

	for _, node := range testNodes {
		testPool.ring.Value = node
		testPool.ring = testPool.ring.Next()
	}

	testPool.complete("localhost:9003")
	expected := nodes{
		{host: "localhost:9000", index: 0, pending: 7},
		{host: "localhost:9002", index: 1, pending: 0},
		{host: "localhost:9001", index: 2, pending: 2},
		{host: "localhost:9003", index: 3, pending: 2},
		{host: "localhost:9004", index: 4, pending: 1},
	}
	require.Equal(t, expected, testPool.nodes)
}

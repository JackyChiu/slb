package slb

import (
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

func TestNodes_newNodes(t *testing.T) {
	testNodes := newNodes(hosts)
	expected := nodes{
		{host: "localhost:9000", index: 0},
		{host: "localhost:9001", index: 1},
		{host: "localhost:9002", index: 2},
		{host: "localhost:9003", index: 3},
		{host: "localhost:9004", index: 4},
	}
	require.Equal(t, expected, testNodes)
}

func TestNodes_Len(t *testing.T) {
	testNodes := newNodes(hosts)
	require.Equal(t, len(hosts), testNodes.Len())
}

func TestNodes_Swap(t *testing.T) {
	testNodes := newNodes(hosts)
	testNodes.Swap(1, 2)
	expected := nodes{
		{host: "localhost:9000", index: 0},
		{host: "localhost:9002", index: 1},
		{host: "localhost:9001", index: 2},
		{host: "localhost:9003", index: 3},
		{host: "localhost:9004", index: 4},
	}
	require.Equal(t, expected, testNodes)
}

func TestNodes_Pop(t *testing.T) {
	testNodes := newNodes(hosts)

	actual := testNodes.Pop()
	expected := &node{
		host:    "localhost:9004",
		pending: 0,
		index:   4,
	}
	require.Equal(t, expected, actual)
}

func TestNodes_Push(t *testing.T) {
	testNodes := newNodes(hosts)

	testNodes.Push(&node{
		host:    "new_host",
		pending: 2,
	})

	expected := nodes{
		{host: "localhost:9000", index: 0, pending: 0},
		{host: "localhost:9001", index: 1, pending: 0},
		{host: "localhost:9002", index: 2, pending: 0},
		{host: "localhost:9003", index: 3, pending: 0},
		{host: "localhost:9004", index: 4, pending: 0},
		{host: "new_host", index: 5, pending: 2},
	}
	require.Equal(t, expected, testNodes)
}

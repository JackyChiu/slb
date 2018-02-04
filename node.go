package slb

import "math"

type node struct {
	host    string
	index   int
	pending int
}

type nodes []*node

func newNode(url string, index int) *node {
	return &node{
		host:    url,
		pending: 0,
		index:   index,
	}
}

func mean(n nodes) (mean float64) {
	length := float64(len(n))
	for _, server := range n {
		mean += float64(server.pending)
	}
	return mean / length
}

func standardDeviation(n nodes) (stdDev float64) {
	length := float64(len(n))
	mean := mean(n)
	for _, server := range n {
		stdDev += math.Pow((float64(server.pending) - mean), 2)
	}
	return math.Sqrt((1 / length) * stdDev)
}

package balance

import "container/heap"

type worker struct {
	host    string
	index   int
	pending int
}

func newWorker(url string, index int) *worker {
	return &worker{
		host:    url,
		pending: 0,
		index:   index,
	}
}

type pool []*worker

func newPool(urls []string) pool {
	var pool pool
	for i, url := range urls {
		worker := newWorker(url, i)
		pool = append(pool, worker)
	}
	heap.Init(pool)
	return pool
}

func (p pool) Len() int {
	return len(p)
}

func (p pool) Less(i, j int) bool {
	return p[i].pending < p[j].pending
}

func (p pool) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
	p[i].index = j
	p[j].index = i
}

func (p pool) Push(x interface{}) {
	w := x.(*worker)
	w.index = p.Len()
	p = append(p, w)
}

func (p pool) Pop() interface{} {
	return p[0]
}

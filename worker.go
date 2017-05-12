package workerpool

import (
	"sync"
	"runtime"
)

type Worker interface {
	Run()
}

type Pool struct {
	tasks chan Worker
	wg    sync.WaitGroup
}

func NewDefault() *Pool {
	numCPUs := runtime.NumCPU()
	return New(numCPUs * 2)
}

func New(maxSize int) *Pool {
	pool := Pool{
		tasks: make(chan Worker),
	}

	pool.wg.Add(maxSize)
	for i := 0; i < maxSize; i++ {
		go func() {
			for t := range pool.tasks {
				t.Run()
			}
			pool.wg.Done()
		}()
	}

	return &pool
}

func (pool *Pool) Submit(w Worker) {
	pool.tasks <- w
}

func (pool *Pool) Shutdown() {
	close(pool.tasks)
	pool.wg.Wait()
}

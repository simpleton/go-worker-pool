package workerpool

import (
	"sync"
	"runtime"
	"time"
	"errors"
)

type Worker interface {
	Run()
}

type Pool struct {
	tasks chan Worker
	sizeWg    sync.WaitGroup
	waitTaskWg sync.WaitGroup
	waitTaskNum int
}

func NewDefault(waitNum int) *Pool {
	numCPUs := runtime.NumCPU()
	if (numCPUs < 4) {
		return New(numCPUs * 2, waitNum)
	} else {
		return New(numCPUs, waitNum)
	}
}

func New(maxSize int, waitNum int) *Pool {
	pool := Pool{
		tasks: make(chan Worker),
		waitTaskNum: waitNum,
	}
	if (waitNum > 0) {
		pool.waitTaskWg.Add(waitNum)
	}
	pool.sizeWg.Add(maxSize)
	for i := 0; i < maxSize; i++ {
		go func() {
			for t := range pool.tasks {
				t.Run()
			}
			pool.sizeWg.Done()
		}()
	}

	return &pool
}

func (pool *Pool) Submit(w Worker) {
	pool.tasks <- w
	if (pool.waitTaskNum > 0) {
		pool.waitTaskWg.Done()
	}
}

func (pool *Pool) Shutdown() (err error) {
	if (pool.waitTaskNum > 0) {
		if waitTimeout(&pool.waitTaskWg, 10 * time.Second) {
			err =  errors.New("Timed out waiting for sync.WaitGroup")
		}
	}
	close(pool.tasks)
	pool.sizeWg.Wait()
	return
}

// waitTimeout waits for the waitgroup for the specified max timeout.
// Returns true if waiting timed out.
func waitTimeout(wg *sync.WaitGroup, timeout time.Duration) bool {
	c := make(chan struct{})
	go func() {
		defer close(c)
		wg.Wait()
	}()
	select {
	case <-c:
		return false // completed normally
	case <-time.After(timeout):
		return true // timed out
	}
}

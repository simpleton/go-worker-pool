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
	wg    sync.WaitGroup
	waitTaskWg *sync.WaitGroup
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
	}
	pool.setWaitTaskNum(waitNum)
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

func (pool *Pool) setWaitTaskNum(size int) {
	if (size > 0) {
		pool.waitTaskWg = new(sync.WaitGroup)
		pool.waitTaskWg.Add(size)
	}
}

func (pool *Pool) Submit(w Worker) {
	if (pool.waitTaskWg != nil) {
		pool.waitTaskWg.Done()
	}
	pool.tasks <- w
}

func (pool *Pool) Shutdown() (err error) {
	if (pool.waitTaskWg != nil) {
		if waitTimeout(pool.waitTaskWg, 3 * time.Second) {
			err =  errors.New("Timed out waiting for wait group")
		}
	}
	close(pool.tasks)
	pool.wg.Wait()
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

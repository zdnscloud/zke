package workerpool

import (
	"errors"
	"sync"
	"sync/atomic"
	"time"
)

var (
	ErrWorkerPoolIsStopped = errors.New("worker pool is stopped")
	ErrWorkerPoolIsTooBusy = errors.New("worker pool is too busy")
)

const SendTaskTimeOut = 3 * time.Second

type Task interface{}

type TaskHandler func(Task)

type WorkerPool struct {
	workerGroup sync.WaitGroup
	handler     TaskHandler
	taskChan    chan Task
	stopped     int32
}

func NewWorkerPool(handler TaskHandler, workerCount int) *WorkerPool {
	if workerCount <= 0 {
		panic("worker count is invalid")
	}

	pool := &WorkerPool{
		handler:  handler,
		taskChan: make(chan Task),
	}
	for i := 0; i < workerCount; i++ {
		pool.workerGroup.Add(1)
		go pool.worker()
	}
	return pool
}

func (pool *WorkerPool) worker() {
	for {
		task, ok := <-pool.taskChan
		if ok == false {
			break
		}
		pool.handler(task)
	}

	pool.workerGroup.Done()
}

func (pool *WorkerPool) SendTask(t Task) error {
	if atomic.LoadInt32(&pool.stopped) == 1 {
		return ErrWorkerPoolIsStopped
	}

	select {
	case pool.taskChan <- t:
		return nil
	case <-time.After(SendTaskTimeOut):
		return ErrWorkerPoolIsTooBusy
	}
}

func (pool *WorkerPool) Stop() error {
	if atomic.CompareAndSwapInt32(&pool.stopped, 0, 1) == false {
		return ErrWorkerPoolIsStopped
	}
	close(pool.taskChan)
	pool.workerGroup.Wait()
	return nil
}

func (pool *WorkerPool) IsStopped() bool {
	return atomic.LoadInt32(&pool.stopped) == 1
}

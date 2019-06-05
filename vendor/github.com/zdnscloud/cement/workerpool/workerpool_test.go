package workerpool

import (
	"sync"
	"testing"
	"time"

	ut "github.com/zdnscloud/cement/unittest"
	"sync/atomic"
)

const SleepTime = 50 * time.Millisecond

func TestWorkerPool(t *testing.T) {
	var initNum int32 = 0
	startTime := time.Now()
	taskCount := 2000
	workerCount := 10
	testComplicateTask(t, &initNum, workerCount, taskCount)
	ut.Equal(t, initNum, int32(taskCount))

	executeTotalInSeq := time.Duration(taskCount) * SleepTime
	executeTotalInParal := executeTotalInSeq / time.Duration(workerCount)
	executeDuration := time.Since(startTime)
	ut.Assert(t, executeDuration < executeTotalInSeq, "actual time should small than %v ", executeTotalInSeq.Seconds())
	ut.Assert(t, executeDuration > executeTotalInParal, "actual time should bigger than %v ", executeTotalInParal.Seconds())
}

func testComplicateTask(t *testing.T, initNum *int32, workerCount, taskCount int) {
	var wg sync.WaitGroup
	worker := func(task Task) {
		i := task.(int32)
		atomic.AddInt32(initNum, i)
		<-time.After(SleepTime)
		wg.Done()
	}
	pool := NewWorkerPool(worker, workerCount)

	for i := 0; i < taskCount; i++ {
		wg.Add(1)
		err := pool.SendTask(int32(1))
		ut.Assert(t, err == nil, "send task get err %v", err)
	}

	wg.Wait()
	ut.Assert(t, pool.Stop() == nil, "pool stop should succeed")
	ut.Assert(t, pool.IsStopped(), "pool should be stopped")
	ut.Assert(t, pool.Stop() == ErrWorkerPoolIsStopped, "pool already stopped")
	ut.Assert(t, pool.Stop() == ErrWorkerPoolIsStopped, "pool already stopped")
}

func TestBusyPool(t *testing.T) {
	taskRunning := make(chan struct{})
	worker := func(task Task) {
		taskRunning <- struct{}{}
		<-time.After(4 * time.Second)
	}
	pool := NewWorkerPool(worker, 1)
	ut.Assert(t, pool.IsStopped() == false, "initially pool isn't stopped")

	err := pool.SendTask(nil)
	ut.Assert(t, err == nil, "send task get err %v", err)
	<-taskRunning
	err = pool.SendTask(nil)
	ut.Assert(t, err == ErrWorkerPoolIsTooBusy, "pool should be busy but %v", err)
	pool.Stop()

	err = pool.SendTask(nil)
	ut.Assert(t, err == ErrWorkerPoolIsStopped, "pool is stop but %v", err)

	ut.Assert(t, pool.Stop() == ErrWorkerPoolIsStopped, "pool already stopped")
}

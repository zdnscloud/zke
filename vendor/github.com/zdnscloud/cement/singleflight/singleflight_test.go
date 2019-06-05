package singleflight

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestDo(t *testing.T) {
	g := New(10)
	v, err := g.Do(uint64(10), func() (interface{}, error) {
		return "bar", nil
	})
	if got, want := fmt.Sprintf("%v (%T)", v, v), "bar (string)"; got != want {
		t.Errorf("Do = %v; want %v", got, want)
	}
	if err != nil {
		t.Errorf("Do error = %v", err)
	}
}

func TestDoErr(t *testing.T) {
	g := New(10)
	someErr := errors.New("Some error")
	v, err := g.Do(uint64(10), func() (interface{}, error) {
		return nil, someErr
	})
	if err != someErr {
		t.Errorf("Do error = %v; want someErr", err)
	}
	if v != nil {
		t.Errorf("unexpected non-nil value %#v", v)
	}
}

func TestDoDupSuppress(t *testing.T) {
	g := New(20)
	c := make(chan string)
	var calls int32
	fn := func() (interface{}, error) {
		atomic.AddInt32(&calls, 1)
		return <-c, nil
	}

	key := uint64(10)
	const n = 10
	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			v, err := g.Do(key, fn)
			if err != nil {
				t.Errorf("Do error: %v", err)
			}
			if v.(string) != "bar" {
				t.Errorf("got %q; want %q", v, "bar")
			}
			wg.Done()
		}()
	}
	time.Sleep(100 * time.Millisecond) // let goroutines above block
	c <- "bar"
	wg.Wait()
	if got := atomic.LoadInt32(&calls); got != 1 {
		t.Errorf("number of calls = %d; want 1", got)
	}
}

func TestInFlightLimit(t *testing.T) {
	g := New(5)
	inDoChan := make(chan int)
	outDoChan := make(chan int)
	fn := func() (interface{}, error) {
		inDoChan <- 1
		<-outDoChan
		return 0, nil
	}

	for i := 0; i < 5; i++ {
		go g.Do(uint64(i), fn)
	}

	for i := 0; i < 5; i++ {
		<-inDoChan
	}

	_, err := g.Do(6, fn)
	if err == nil {
		t.Errorf("limit exceed should return err")
	}

	if g.inflightCount != 5 {
		t.Errorf("should have 5 inflight do but get %d", g.inflightCount)
	}

	//make two do return
	for i := 0; i < 2; i++ {
		outDoChan <- 1
	}

	for g.inflightCount != 3 {
		<-time.After(1 * time.Second)
	}

	_, err = g.Do(7, func() (interface{}, error) {
		return 10, nil
	})
	if err != nil {
		t.Errorf("limit isn't exceeded no err should be returned")
	}
}

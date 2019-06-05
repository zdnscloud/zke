package singleflight

import (
	"errors"
	"sync"
)

var ErrExceedLimit = errors.New("exceed the concurrent limit")

type call struct {
	wg  sync.WaitGroup
	val interface{}
	err error
}

type Group struct {
	mu              sync.Mutex       // protects m
	m               map[uint64]*call // lazily initialized
	inflightCount   uint32
	concurrentLimit uint32
}

func New(concurrentLimit uint32) *Group {
	return &Group{
		m:               make(map[uint64]*call),
		inflightCount:   0,
		concurrentLimit: concurrentLimit,
	}
}

func (g *Group) Do(key uint64, fn func() (interface{}, error)) (interface{}, error) {
	g.mu.Lock()
	if c, ok := g.m[key]; ok {
		g.mu.Unlock()
		c.wg.Wait()
		return c.val, c.err
	}

	if g.concurrentLimit > 0 && g.inflightCount+1 > g.concurrentLimit {
		g.mu.Unlock()
		return nil, ErrExceedLimit
	}

	c := new(call)
	c.wg.Add(1)
	g.m[key] = c
	g.inflightCount += 1
	g.mu.Unlock()

	c.val, c.err = fn()
	c.wg.Done()

	g.mu.Lock()
	if g.concurrentLimit > 0 {
		g.inflightCount -= 1
	}
	delete(g.m, key)
	g.mu.Unlock()

	return c.val, c.err
}

package rcu

import (
	"sync"
	"sync/atomic"
	"time"
)

type RCU[T any] struct {
	value  *atomic.Pointer[T]
	update func(*T) *T
	ticker *time.Ticker
	done   chan struct{}
	mu     *sync.Mutex
}

func NewRCU[T any](updateInterval time.Duration, update func(*T) *T) *RCU[T] {
	self := &RCU[T]{
		value:  &atomic.Pointer[T]{},
		update: update,
		ticker: time.NewTicker(updateInterval),
		done:   make(chan struct{}),
		mu:     &sync.Mutex{},
	}
	self.forceUpdate()
	go self.scheduleUpdate()
	return self
}

func (r *RCU[T]) Load() *T {
	return r.value.Load()
}

func (r *RCU[T]) scheduleUpdate() {
	for {
		select {
		case <-r.done:
			return
		case <-r.ticker.C:
			r.forceUpdate()
		}
	}
}

func (r *RCU[T]) forceUpdate() {
	r.mu.Lock()
	defer r.mu.Unlock()

	newValue := r.update(r.value.Load())
	r.value.Store(newValue)
}

func (r *RCU[T]) Close() {
	r.ticker.Stop()
	r.done <- struct{}{}
}

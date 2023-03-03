package rxxr

import (
	"sync"
)

type Pipe[T any] struct {
	panicHandler func(any)

	sendOnSubscribe bool
	value           *T

	lock      sync.RWMutex
	listeners map[chan T]func()
	closeOnce sync.Once
	publish   chan T
}

// New creates a simple pipe that transmits received values to all subscribers
func New[T any](options ...Option[T]) *Pipe[T] {
	p := &Pipe[T]{
		listeners: map[chan T]func(){},
		publish:   make(chan T),
	}
	go p.send()

	for _, o := range options {
		if o != nil {
			o(p)
		}
	}

	return p
}

func (p *Pipe[T]) Subscribe(fn func(T)) (cancel func()) {
	if fn != nil {
		p.lock.Lock()
		defer p.lock.Unlock()
		if p.publish == nil {
			return func() {}
		}

		c := make(chan T)
		closeC := makeCloser(c)
		p.listeners[c] = closeC
		go runLoop(fn, p.panicHandler, onSubscribe(fn, p.panicHandler, p.value, p.sendOnSubscribe), c)
		return p.unsubscribe(closeC, c)
	}
	return func() {}
}

func (p *Pipe[T]) Publish(t T) {
	p.lock.Lock()
	defer p.lock.Unlock()
	if p.publish != nil {
		p.value = &t
		p.publish <- t
	}
}

func (p *Pipe[T]) Value() (t T, ok bool) {
	p.lock.RLock()
	defer p.lock.RUnlock()
	if p.value != nil {
		return *p.value, true
	}
	return
}

func (p *Pipe[T]) Close() {
	p.closeOnce.Do(func() {
		p.lock.Lock()
		defer p.lock.Unlock()
		close(p.publish)
		p.publish = nil
	})
}

func (p *Pipe[T]) send() {
	defer p.closeAll()
	for v := range p.publish {
		p.sendAll(v)()
	}
}

func (p *Pipe[T]) sendAll(v T) func() {
	var wg sync.WaitGroup
	p.lock.RLock()
	defer p.lock.RUnlock()
	for c := range p.listeners {
		wg.Add(1)
		go func(c chan T) {
			defer wg.Done()
			c <- v
		}(c)
	}
	return wg.Wait
}

func (p *Pipe[T]) closeAll() {
	p.lock.Lock()
	defer p.lock.Unlock()
	for _, v := range p.listeners {
		v()
	}
	p.listeners = nil
}

func makeCloser[T any](c chan T) func() {
	var once sync.Once
	return func() {
		once.Do(func() { close(c) })
	}
}

func (p *Pipe[T]) unsubscribe(closeC func(), c chan T) func() {
	return func() {
		p.lock.Lock()
		defer p.lock.Unlock()
		delete(p.listeners, c)
		closeC()
	}
}

func onSubscribe[T any](fn func(T), panicHandler func(any), value *T, sendOnSubscribe bool) func() {
	if sendOnSubscribe && value != nil {
		return func() {
			defer maybeHandlePanic(panicHandler)
			fn(*value)
		}
	}
	return func() {}
}

func runLoop[T any](fn func(T), panicHandler func(any), sendOnSubscribe func(), c chan T) {
	sendOnSubscribe()
	for v := range c {
		func() {
			defer maybeHandlePanic(panicHandler)
			fn(v)
		}()
	}
}

func maybeHandlePanic(panicHandler func(any)) {
	if panicHandler != nil {
		if r := recover(); r != nil {
			panicHandler(r)
		}
	}
}

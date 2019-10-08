package pool

import (
	"sync"
)

// Runnable is everything which can be started/stopped
type runnable interface {
	Start()
	Stop()
}

// Pool manages a list of runnables and can Start/Stop them
type Pool struct {
	runnables []runnable

	done []chan struct{}
	wg   sync.WaitGroup
}

// Add a new runnables to the pool
func (p *Pool) Add(r runnable) {
	p.runnables = append(p.runnables, r)
}

// Len returns the number of runnables currently in the pool
func (p *Pool) Len() int {
	return len(p.runnables)
}

// StartAll runs all runnables
func (p *Pool) StartAll() {
	for _, r := range p.runnables {
		go func(r runnable) {
			p.wg.Add(1)
			r.Start()
			defer p.wg.Done()
		}(r)
	}
}

// StopAll stops all runnables
func (p *Pool) StopAll() {
	for _, r := range p.runnables {
		r.Stop()
	}
	p.wg.Wait()
	p.done = nil
}

package utils

import (
	"sync"
)

type Work func()

type WorkerPool struct {
	work chan Work
	wg   sync.WaitGroup
}

func NewWorkerPool(maxGoroutines int) *WorkerPool {
	p := WorkerPool{
		work: make(chan Work),
	}

	p.wg.Add(maxGoroutines)
	for i := 0; i < maxGoroutines; i++ {
		go func() {
			for w := range p.work {
				w()
			}
			p.wg.Done()
		}()
	}

	return &p
}

func (p *WorkerPool) Run(w Work) {
	p.work <- w
}

func (p *WorkerPool) Shutdown() {
	close(p.work)
	p.wg.Wait()
}

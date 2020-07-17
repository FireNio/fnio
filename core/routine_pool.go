package core

import (
	"sync/atomic"
	"unsafe"
)

type RoutinePool struct {
	running int32
	jobs    *ArrayBlockingQueue
	cap     int
}

func NewRoutinePool(queueSize int, routineSize int) *RoutinePool {
	return &RoutinePool{running: 1, jobs: NewArrayBlockingQueue(queueSize), cap: routineSize}
}

func (p *RoutinePool) Start() {
	for i := 0; i < p.cap; i++ {
		go func(goid int64) {
			jobs := p.jobs
			for p.IsRunning() {
				job := jobs.PollTimeout()
				if job != nil {
					(*(*Runnable)(job)).Run(goid)
				}
			}
		}(NextGoId())
	}
}

func (p *RoutinePool) Stop(goid int64) {
	atomic.StoreInt32(&p.running, 0)
	jobs := p.jobs
	for {
		job := jobs.Poll()
		if job == nil {
			break
		}
		(*(*Runnable)(job)).Run(goid)
	}
	p.jobs.WakeUpAll()
}

func (p *RoutinePool) IsRunning() bool {
	return atomic.LoadInt32(&p.running) == 1
}

func (p *RoutinePool) Submit(job Runnable) bool {
	jobs := p.jobs
	var jobp = unsafe.Pointer(&job)
	if !jobs.Offer(jobp) {
		return false
	}
	return !(!p.IsRunning() && jobs.Remove(jobp))
}

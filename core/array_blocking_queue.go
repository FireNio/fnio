package core

import (
	"sync"
	"sync/atomic"
	"unsafe"
)

type ArrayBlockingQueue struct {
	cond       *sync.Cond
	items      []unsafe.Pointer
	put_index  int32
	take_index int32
	size       int32
	cap        int32
	mask       int32
}

func NewArrayBlockingQueue(cap int) *ArrayBlockingQueue {
	if !IsSq2(cap) {
		panic("Not index number")
	}
	var m = sync.Mutex{}
	return &ArrayBlockingQueue{items: make([]unsafe.Pointer, cap),
		cap:  int32(cap),
		mask: int32(cap - 1),
		cond: sync.NewCond(&m),
	}
}

func (q *ArrayBlockingQueue) Offer(item unsafe.Pointer) bool {
	cond := q.cond
	cond.L.Lock()
	defer cond.L.Unlock()
	if q.size == q.cap {
		return false
	}
	q.items[q.put_index] = item
	q.put_index = inc_index(q.put_index, q.mask)
	//atomic.AddInt32(&q.size, 1)
	q.size++
	cond.Signal()
	return true
}

func (q *ArrayBlockingQueue) PollTimeout() unsafe.Pointer {
	cond := q.cond
	cond.L.Lock()
	defer cond.L.Unlock()
	if q.size == 0 {
		cond.Wait()
		if q.size == 0 {
			return nil
		}
	}
	pointer := q.items[q.take_index]
	q.items[q.take_index] = nil
	q.take_index = inc_index(q.take_index, q.mask)
	//atomic.AddInt32(&q.size, -1)
	q.size--
	return pointer
}

func (q *ArrayBlockingQueue) Poll() unsafe.Pointer {
	cond := q.cond
	cond.L.Lock()
	defer cond.L.Unlock()
	if q.size == 0 {
		return nil
	}
	pointer := q.items[q.take_index]
	q.items[q.take_index] = nil
	q.take_index = inc_index(q.take_index, q.mask)
	//atomic.AddInt32(&q.size, -1)
	q.size--
	return pointer
}

func (q *ArrayBlockingQueue) Size() int32 {
	return atomic.LoadInt32(&q.size)
}

func (q *ArrayBlockingQueue) Remove(item unsafe.Pointer) bool {
	cond := q.cond
	cond.L.Lock()
	defer cond.L.Unlock()
	if q.size == 0 {
		return false
	}
	take_index := q.take_index
	put_index := q.put_index
	items := q.items
	mask := q.mask
	for take_index != put_index {
		if item == items[take_index] {
			items[take_index] = nil
			for take_index != put_index {
				var next_take_index = inc_index(take_index, mask)
				items[take_index] = items[next_take_index]
				take_index = next_take_index
			}
			return true
		}
		take_index = inc_index(take_index, mask)
	}
	return false
}

func (q *ArrayBlockingQueue) IsEmpty() bool {
	return q.Size() == 0
}

func (q *ArrayBlockingQueue) Cap() int32 {
	return q.cap
}

func (q *ArrayBlockingQueue) WakeUpAll() {
	cond := q.cond
	cond.L.Lock()
	cond.Broadcast()
	cond.L.Unlock()
}

func inc_index(index int32, mask int32) int32 {
	return (index + 1) & mask
}

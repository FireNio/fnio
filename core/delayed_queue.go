package core

const (
	DELAY_CANCEL_MASK uint64 = 1 << 63
	DEALY_DONE_MASK   uint64 = 1 << 62
	DELAY_DELAY_MASK  uint64 = ^(DELAY_CANCEL_MASK | DEALY_DONE_MASK)
)

type DelayTask struct {
	flags uint64
}

func (d *DelayTask) Cancel() {
	d.flags |= DELAY_CANCEL_MASK
}

func (d *DelayTask) GetDelay() int64 {
	return int64(d.flags & DELAY_DELAY_MASK)
}

func (d *DelayTask) IsCanceled() bool {
	return (d.flags & DELAY_CANCEL_MASK) != 0
}

func (d *DelayTask) IsDone() bool {
	return (d.flags & DEALY_DONE_MASK) != 0
}

func (d *DelayTask) Done() {
	d.flags |= DEALY_DONE_MASK
}

func (d *DelayTask) CompareTo(other DelayRunnable) int64 {
	return d.GetDelay() - other.GetDelay()
}

type DelayedQueue struct {
	queue *[]DelayRunnable
	size  int
}

func NewDelayedQueue(size int) *DelayedQueue {
	var queue = make([]DelayRunnable, size)
	return &DelayedQueue{queue: &queue}
}

func (q *DelayedQueue) Clear() {
	size := q.size
	queue := *q.queue
	for i := 0; i < size; i++ {
		var t = queue[i]
		if t != nil {
			queue[i] = nil
		}
	}
	q.size = 0
}

func (q *DelayedQueue) Contains(x DelayRunnable) bool {
	return q.indexOf(x) != -1
}

func (q *DelayedQueue) finishPoll(f DelayRunnable) DelayRunnable {
	q.size--
	queue := *q.queue
	var s = q.size
	var x = queue[s]
	queue[s] = nil
	if s != 0 {
		q.siftDown(0, x)
	}
	return f
}

func (q *DelayedQueue) grow() {
	queue := *q.queue
	var oldCapacity = len(queue)
	var newCapacity = oldCapacity + (oldCapacity >> 1) // grow 50%
	if newCapacity < 0 {
		// overflow
		panic("overflow")
	}
	var new_queue = make([]DelayRunnable, newCapacity)
	copy(new_queue[:], queue)
	q.queue = &new_queue
}

func (q *DelayedQueue) indexOf(x DelayRunnable) int {
	if x != nil {
		queue := *q.queue
		var size = q.size
		for i := 0; i < size; i++ {
			if x == queue[i] {
				return i
			}
		}
	}
	return -1
}

func (q *DelayedQueue) IsEmpty() bool {
	return q.size == 0
}

func (q *DelayedQueue) Offer(e DelayRunnable) bool {
	var i = q.size
	queue := *q.queue
	if i >= len(queue) {
		q.grow()
	}
	q.size = i + 1
	if i == 0 {
		queue[0] = e
	} else {
		q.siftUp(i, e)
	}
	return true
}

func (q *DelayedQueue) Peek() DelayRunnable {
	return (*q.queue)[0]
}

func (q *DelayedQueue) Poll() DelayRunnable {
	if q.IsEmpty() {
		return nil
	}
	return q.finishPoll((*q.queue)[0])
}

func (q *DelayedQueue) Remove(x DelayRunnable) bool {
	var i = q.indexOf(x)
	if i < 0 {
		return false
	}
	queue := *q.queue
	q.size--
	var s = q.size
	var replacement = queue[s]
	queue[s] = nil
	if s != i {
		q.siftDown(i, replacement)
		if queue[i] == replacement {
			q.siftUp(i, replacement)
		}
	}
	return true
}

func (q *DelayedQueue) siftDown(k int, key DelayRunnable) {
	size := q.size
	var half = q.size >> 1
	queue := *q.queue
	for ; k < half; {
		var child = (k << 1) + 1
		var c = queue[child]
		var right = child + 1
		if right < size && c.CompareTo(queue[right]) > 0 {
			child = right
			c = queue[child]
		}
		if key.CompareTo(c) <= 0 {
			break
		}
		queue[k] = c
		k = child
	}
	queue[k] = key
}

func (q *DelayedQueue) siftUp(k int, key DelayRunnable) {
	queue := *q.queue
	for ; k > 0; {
		var parent = (k - 1) >> 1
		var e = queue[parent]
		if key.CompareTo(e) >= 0 {
			break
		}
		queue[k] = e
		k = parent
	}
	queue[k] = key
}

func (q *DelayedQueue) Size() int {
	return q.size
}

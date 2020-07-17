package core

import (
	"fmt"
	"testing"
)

func TestDelayedQueue(t *testing.T) {

	var q = NewDelayedQueue(4)

	q.Offer(NewTestDelayTask(200))
	q.Offer(NewTestDelayTask(100))
	q.Offer(NewTestDelayTask(400))
	q.Offer(NewTestDelayTask(300))
	q.Offer(NewTestDelayTask(500))
	q.Offer(NewTestDelayTask(600))

	q.Poll().Run(0)
	q.Poll().Run(0)
	q.Poll().Run(0)
	q.Poll().Run(0)
	q.Poll().Run(0)
	q.Poll().Run(0)
}

type TestDelayTask struct {
	DelayTask
}

func (d *TestDelayTask) Run(goid int64) {
	fmt.Println("delay: ", d.DelayTask.GetDelay())
}

func NewTestDelayTask(delay int64) *TestDelayTask {
	var d = TestDelayTask{}
	d.flags = uint64(delay)
	return &d
}

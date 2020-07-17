package core

import (
	"fmt"
	"runtime"
	"testing"
	"time"
)

func TestRoutinePool(t *testing.T) {
	runtime.GOMAXPROCS(4)

	p := NewRoutinePool(16, 2)

	p.Start()

	var r1 = &TestPoolRunnable{msg: "msg1"}
	var r2 = &TestPoolRunnable{msg: "msg2"}
	var r3 = &TestPoolRunnable{msg: "msg3"}
	var r4 = &TestPoolRunnable{msg: "msg4"}
	var r5 = &TestPoolRunnable{msg: "msg5"}

	p.Submit(r1)
	p.Submit(r2)
	p.Submit(r3)
	p.Submit(r4)
	p.Submit(r5)

	p.Stop(0)




}

type TestPoolRunnable struct {
	msg string
}

func (t * TestPoolRunnable) Run(goid int64)  {
	fmt.Println(time.Now().UnixNano(),goid, t.msg)
}


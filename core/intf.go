package core

import (
	"sync/atomic"
	"unsafe"
)

var _goid int64 = 1

func NextGoId() int64 {
	return atomic.AddInt64(&_goid, 1)
}

type Runnable interface {
	Run(goid int64)
}

type DelayRunnable interface {
	Runnable
	Cancel()
	GetDelay() int64
	IsCanceled() bool
	IsDone() bool
	Done()
	CompareTo(other DelayRunnable) int64
}

type Queue interface {
	Offer(item unsafe.Pointer) bool
	Poll() unsafe.Pointer
	IsEmpty() bool
}



package core

type Runnable interface {
	Run()
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

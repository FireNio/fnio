package core

import (
	"fmt"
	"runtime"
	"testing"
	"time"
	"unsafe"
)

func TestArrayBlockingQueue(t *testing.T) {
	num := runtime.NumCPU()
	runtime.GOMAXPROCS(num)

	q := NewArrayBlockingQueue(16)

	var s100 = "100"
	var s200 = "200"
	var s300 = "300"
	var s400 = "400"

	go poll(q)

	time.Sleep(1000 * time.Millisecond)
	fmt.Println("start offer...")
	q.Offer(unsafe.Pointer(&s200))
	q.Offer(unsafe.Pointer(&s100))
	q.Offer(unsafe.Pointer(&s400))
	q.Offer(unsafe.Pointer(&s300))

	q.Remove(unsafe.Pointer(&s100))


	go fmt.Println(*(*string)(q.Poll()))
	go fmt.Println(*(*string)(q.Poll()))



}

func poll(q * ArrayBlockingQueue)  {
	fmt.Println("start poll...")
	fmt.Println(*(*string)(q.Poll()))
}

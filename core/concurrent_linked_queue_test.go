package core

import (
	"fmt"
	"testing"
	"unsafe"
)

func TestQueue(t *testing.T) {
	queue := NewConcurrentLinkedQueue()

	var s = "abc"

	var ands = &s

	var vs = unsafe.Pointer(ands)

	queue.Offer(vs)

	value := (*string)(queue.Poll())

	var eq = value == ands

	fmt.Println(eq)

	fmt.Println(*value)
}

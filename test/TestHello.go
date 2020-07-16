package main

import (
	"fmt"
	"github.com/firenio/firenio-go/core"
	"unsafe"
)

func main() {

	fmt.Println("123")

	TestIntPtrMap()
}

func TestIntPtrMap()  {

	ptrMap := core.NewIntPtrMap(16, 0.5)

	var  s = "abc"

	var ands = &s

	var vs = unsafe.Pointer(ands)

	ptrMap.Put(1, vs)

	ptrMap.Put(33, vs)

	ptrMap.Remove(1)

	get := (*string)(ptrMap.Get(33))

	fmt.Println(*get)

}

func TestQueue()  {
	queue := core.NewConcurrentLinkedQueue()

	var  s = "abc"

	var ands = &s

	var vs = unsafe.Pointer(ands)

	queue.Offer(vs)

	value := (*string) (queue.Poll())

	var eq = value == ands

	fmt.Println(eq)

	fmt.Println(*value)
}

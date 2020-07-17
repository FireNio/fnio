package core

import (
	"fmt"
	"testing"
	"unsafe"
)

func TestIntPtrMap(t *testing.T) {

	ptrMap := NewIntPtrMap(16, 0.5)

	var  s = "abc"

	var ands = &s

	var vs = unsafe.Pointer(ands)

	ptrMap.Put(1, vs)

	ptrMap.Put(33, vs)

	ptrMap.Remove(1)

	get := (*string)(ptrMap.Get(33))

	fmt.Println(*get)

}
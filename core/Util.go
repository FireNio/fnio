package core

import (
	"unsafe"
)

func MaxInt32(a, b int32) int32 {
	if a > b {
		return a
	} else {
		return b
	}
}

func ClothCoverInt32(v int32) int32 {
	var n int32 = 2
	for ; n < v; {
		n <<= 1
	}
	return n
}

func FillInt32(array []int32, value int32) {
	var l = len(array)
	for i := 0; i < l; i++ {
		array[i] = value
	}
}

func FillInt64(array []int64, value int64) {
	var l = len(array)
	for i := 0; i < l; i++ {
		array[i] = value
	}
}

func FillPtrNil(array []unsafe.Pointer) {
	var l = len(array)
	for i := 0; i < l; i++ {
		array[i] = nil
	}
}

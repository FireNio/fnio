package main

import (
	"fmt"
	"strconv"
)

func main() {

	fmt.Println("123")

	var CANCEL_MASK uint8 = 1 << 7
	var DELAY_MASK  uint8 = ^(CANCEL_MASK)

	fmt.Println(strconv.FormatInt(int64(DELAY_MASK), 2))
}

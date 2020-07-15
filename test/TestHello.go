package main

import (
	"fmt"
	"github.com/firenio/firenio-go/core"
)

func main() {

	fmt.Println("123")

	core.Run()

	queue := core.NewConcurrentLinkedQueue()

	var  v = "abc"

	queue.Offer(&v)

	value := queue.Poll()

	fmt.Println(value)
}



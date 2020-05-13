package main

import (
	"fmt"
	"time"
)

//Hello, Harry!
//todo 一定要结合MPG模型，调度器原理考虑问题额
func main() {
	name := "Eric"
	//--
	go func() {
		fmt.Printf("Hello, %s!\n", name)
	}()
	//--
	name = "Harry"
	time.Sleep(time.Millisecond)
	// time.Sleep(time.Millisecond)
	// name = "Harry"
}

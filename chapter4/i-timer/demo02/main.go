package main

import (
	"fmt"
	"time"
)

//使用timer定时器，超时后需要重置，才能继续触发。
func main() {

	d := time.Duration(time.Second*2)
	t := time.NewTimer(d)
	defer t.Stop()
	for {
		<- t.C
		fmt.Println("timer超时...")
		// need reset
		t.Reset(time.Second*2)
	}
}

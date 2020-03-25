package main

import "runtime"

func main() {
	go println("Go! Goroutine!")

	//方式一  此种不一定是100%的，当调度复杂的时候
	//time.Sleep(time.Millisecond)
	//方式二 用runtime.Gosched()，暂停当前的G
	runtime.Gosched()
}

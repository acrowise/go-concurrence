package main

import (
	"fmt"
	"time"
)
func main() {

	fmt.Println(time.Now())
    pause:=10 * time.Second
	//创建一个每隔10s执行回调函数的一个定时器
	timer:=time.AfterFunc(pause, func() {
		fmt.Println(time.Now())
	} )
	//每隔12s重置定时器
	for{
		time.Sleep(12 * time.Second)
		timer.Reset(pause)
	}
}

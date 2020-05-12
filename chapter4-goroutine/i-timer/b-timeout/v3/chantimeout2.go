package main

import (
	"fmt"
	"time"
)


//select语句与for语句连用可以持续地从一个通道接收元素值。
//但是，若每次接收时都初始化一个定时器显然有些浪费，好在定时器是可以复用的。(方式二)


func main() {


	intChan := make(chan int, 1)
	//发送
	go func() {
		for i := 0; i < 5; i++ {
			time.Sleep(3 * time.Second)
			intChan <- i
		}
		close(intChan)
	}()
	//接收
	timeout := 2 * time.Second
	timer:=time.NewTimer(timeout)

	for {
		select {
		case e, ok := <-intChan:
			if !ok {
				fmt.Println("End.")
				return
			}
			fmt.Printf("Received: %v\n", e)
		case <-timer.C:
			fmt.Println("超时啦，生产端能不能提速啊!")
		}
		//todo 定时器到期之后到重置之前不可以再接收了
		//<-timer.C
		timer.Reset(timeout)
	}





}

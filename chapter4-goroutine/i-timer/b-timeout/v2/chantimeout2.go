package main

import (
	"fmt"
	"time"
)


//select语句与for语句连用可以持续地从一个通道接收元素值。
//但是，若每次接收时都初始化一个定时器显然有些浪费，好在定时器是可以复用的。(方式一)

func main() {


	intChan := make(chan int, 1)

	//每隔3秒发送一个元素，发送完5个关闭通道
	go func() {
		for i := 0; i < 5; i++ {
			time.Sleep(3 * time.Second)
			intChan <- i
		}
		close(intChan)
	}()


   //接收
	timeout := 2 * time.Second
	var timer *time.Timer
	for {
		if timer == nil {
			timer = time.NewTimer(timeout)
		} else {
			timer.Reset(timeout)
		}
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
	}





}

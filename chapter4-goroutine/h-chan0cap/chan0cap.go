package main

import (
	"fmt"
	"time"
)

func main() {

	//故意让发送间隔为1s接收间隔为2s，但最终你会发现间隔都是2s

	sendingInterval := time.Second
	receptionInterval := time.Second * 2
	intChan := make(chan int, 0)

	//发送操作-------------------------------------
	go func() {
		var ts0, ts1 int64
		for i := 1; i <= 5; i++ {
			intChan <- i //往通道中发送元素
			ts1 = time.Now().Unix()//记录发送结束的时间
			if ts0 == 0 {//为0表示发送的是第一个元素
				fmt.Println("Sent:", i)
			} else {
				//统计一下发送了什么，以及发送成功之后的间隔时间
				fmt.Printf("Sent: %d [interval: %d s]\n", i, ts1-ts0)
			}
			ts0 = time.Now().Unix() //下一次发送元素前的初始时间
			time.Sleep(sendingInterval)
		}
		close(intChan)
	}()
	var ts0, ts1 int64

//接收操作----------------------------------------------------
Loop:
	for {
		select {
		case v, ok := <-intChan:
			//如果管道关闭了就跳出来
			if !ok {
				break Loop
			}
			ts1 = time.Now().Unix()
			if ts0 == 0 {
				fmt.Println("Received:", v)
			} else {
				fmt.Printf("Received: %d [interval: %d s]\n", v, ts1-ts0)
			}
		}
		ts0 = time.Now().Unix()
		time.Sleep(receptionInterval)
	}
	fmt.Println("End.")
}

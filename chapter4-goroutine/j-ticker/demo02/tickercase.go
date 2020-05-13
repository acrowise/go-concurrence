package main

import (
	"fmt"
	"time"
)

func main() {
	//创建一个通信用的通道
	intChan := make(chan int, 1)
	//创建一个断续器
	ticker := time.NewTicker(time.Second)



	//发送方
	go func() {
		for _ = range ticker.C {
			select {
			case intChan <- 1:
			case intChan <- 2:
			case intChan <- 3:
			}
		}
		//考虑一下为什么总是打印不出这一句呢，因为for一直没有执行完毕啊
		fmt.Println("End. [sender]")
	}()

    //接收方
	for e := range intChan {
		fmt.Printf("Received: %v\n", e)
	}
}

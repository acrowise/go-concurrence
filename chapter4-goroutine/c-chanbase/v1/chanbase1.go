package main

import (
	"fmt"
	"time"
)

var strChan = make(chan string, 3)
//演示的是对strChan的存取操作
//syncChan1与syncChan2是作为辅助通信来用的

func main() {
	syncChan1 := make(chan struct{}, 1) //标志可以接收了
	syncChan2 := make(chan struct{}, 2) //记录发送以及接收2G都结束了

	//接收操作----------------------------
	go func() {
		<-syncChan1
		fmt.Println("Received a sync signal and wait a second... [receiver]")
		time.Sleep(time.Second)
		for {
			if elem, ok := <-strChan; ok {
				fmt.Println("Received:", elem, "[receiver]")
			} else {
				break
			}
		}
		fmt.Println("Stopped. [receiver]")
		syncChan2 <- struct{}{}
	}()
	//发送操作-------------------------------
	go func() {
		for _, elem := range []string{"a", "b", "c", "d"} {
			strChan <- elem
			fmt.Println("Sent:", elem, "[sender]")
			if elem == "c" {
				//当发送到c的时候，发个信号让接收方开始接收吧
				syncChan1 <- struct{}{}
				fmt.Println("Sent a sync signal. [sender]")
			}
		}
		fmt.Println("Wait 2 seconds... [sender]")
		time.Sleep(time.Second * 2)
		close(strChan)
		syncChan2 <- struct{}{}
	}()
	//阻塞
	<-syncChan2
	<-syncChan2
}

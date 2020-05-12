package main

import (
	"fmt"
	"time"
)

//todo 注意:使用定时器，我们可以便捷地实现对接收操作的超时设定,要深刻体会额
func main() {
	intChan := make(chan int, 1)


	//发送
	go func() {
		time.Sleep(10 *time.Second)
		intChan <- 1
	}()
	//select {
	//case e := <-intChan:
	//	fmt.Printf("Received: %v\n", e)
	//case <-time.NewTimer(5 *time.Second).C:
	//	fmt.Println("5s超时额~")
	//}

   //接收
	select {
	case e := <-intChan:
		fmt.Printf("Received: %v\n", e)
	case <-time.After(5 *time.Second):
		fmt.Println("5s超时额~")
	}



}

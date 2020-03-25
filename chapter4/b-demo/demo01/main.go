package main

import (
	"fmt"
	"time"
)


//这种所有线程或者进程(应该指的指的是除了g0之外的所有g吧)都在等待资源释放的情况，我们便把它称之为死锁。  fatal error: all goroutines are asleep - deadlock!
//
//死锁是一个非常有意思的话题，常见的死锁大致分为以下几类：
//i. 只在单一goroutine里操作信道，例子如上。
//ii. 串联信道中间一环挂起



func main() {


	var  strChan chan  string
	strChan = make(chan string,3)

	elem:=<-strChan

	fmt.Println(elem,"hello")

	time.Sleep(1 * time.Second)

	fmt.Print("hi")
	time.Sleep(1 * time.Second)
	fmt.Print("hi2")

}

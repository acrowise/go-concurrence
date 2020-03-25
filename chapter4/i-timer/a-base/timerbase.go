package main

import (
	"fmt"
	"time"
)

func main() {

	//返回 *time.Timer类型    timer.Reset()    timer.Stop()

	fmt.Printf("初始化时的绝对时间: %v.\n", time.Now())
	interval:=10 * time.Second
	timer := time.NewTimer(interval)
	fmt.Printf("相对到期时间: %v.\n", interval)

	//会一直阻塞，直到定时器到期
	expirationTime := <-timer.C
	fmt.Printf("绝对到期时间为: %v.\n", expirationTime)

	//这个时候调用停止定时器的结果是false,因为定时器在这个时候已经过期了
	fmt.Printf("Stop timer: %v.\n", timer.Stop())
}

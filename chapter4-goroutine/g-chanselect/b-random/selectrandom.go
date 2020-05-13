package main

import "fmt"

//@todo 这里未死锁的原因是始终未用满缓冲通道的值
func main() {
	chanCap := 5
	intChan := make(chan int, chanCap)


	for i := 0; i < chanCap; i++ {
		select {
		case intChan <- 1:
		case intChan <- 2:
		case intChan <- 3:
		}
	}



	for i := 0; i < chanCap; i++ {
		fmt.Printf("%d\n", <-intChan)
	}
}

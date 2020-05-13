package main

import "fmt"

func main() {
	//intChan
	intChan := make(chan int, 10)
	for i := 0; i < 10; i++ {
		intChan <- i
	}
	close(intChan)

	syncChan := make(chan struct{}, 1)
	//todo 当 for-select组合，则break只能结束select，for仍然继续死循环下去，所以得利用标签位置额
	go func() {
	Loop:
		for {
			select {
			case e, ok := <-intChan:
				if !ok {
					fmt.Println("End.")
					//break
					break Loop
				}
				fmt.Printf("Received: %v\n", e)
			}
		}
		syncChan <- struct{}{}
	}()
	<-syncChan
}

package main

import (
	"fmt"
	"time"
)

var mapChan = make(chan map[string]int, 1)
//todo 演示通道元素值的传递的是副本还是什么，是否会相互影响
func main() {
	syncChan := make(chan struct{}, 2)
	// 用于演示接收操作------------------------
	go func() {
		for {
			if elem, ok := <-mapChan; ok {
				elem["count"]++
			} else {
				break
			}
		}
		fmt.Println("Stopped. [receiver]")
		syncChan <- struct{}{}
	}()
	// 用于演示发送操作-------------------
	go func() {
		countMap := make(map[string]int)
		for i := 0; i < 5; i++ {
			mapChan <- countMap
			//之所以等一会打印，是想让接收方先把做修改，看看是否会影响发送方中的元素
			time.Sleep(time.Millisecond)
			//todo 对countMap中的值的修改却是在接收方，但是由于map是引用类型，所以在发送方打印的时候就看到了端倪
			fmt.Printf("The count map: %v. [sender]\n", countMap)
		}
		close(mapChan)
		syncChan <- struct{}{}
	}()



	<-syncChan
	<-syncChan
}

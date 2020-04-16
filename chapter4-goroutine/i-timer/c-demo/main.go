package main

import (
	"fmt"
	"time"
)

func main() {

	fmt.Println(time.Now())
    pause:=10 * time.Second

	//异步执行的
	timer:=time.AfterFunc(pause, func() {

		fmt.Println("好苦啊")
		fmt.Println(time.Now())
	} )



	fmt.Println("开始...")
	for{
		fmt.Println("等待...")
		time.Sleep(12 * time.Second)
		timer.Reset(pause)
	}




}

package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	//声明一个读写锁变量
	var rwm sync.RWMutex

    //并发开启g用于对读写锁rwm的读锁定和读解锁操作
	for i := 0; i < 3; i++ {
		go func(i int) {
			//g尝试上读锁
			fmt.Printf("Try to lock for reading... [%d]\n", i)
			rwm.RLock()
			fmt.Printf("Locked for reading. [%d]\n", i)
			//g故意休眠2s观察数据
			time.Sleep(time.Second * 2)
			fmt.Printf("Try to unlock for reading... [%d]\n", i)
			//g解开读锁
			rwm.RUnlock()
			fmt.Printf("Unlocked for reading. [%d]\n", i)
		}(i)
	}
	//main故意睡一会，让g们运行起来
	time.Sleep(time.Millisecond * 100)



	//main尝试上写锁，此处肯定会阻塞，因为必须等所有的g都读解锁完成
	fmt.Println("Try to lock for writing...")
	rwm.Lock()
	fmt.Println("Locked for writing.")



}

package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	//初始化一个互斥锁
	var mutex sync.Mutex

	//先让main将锁锁起来
	fmt.Println("Lock the lock. (main)")
	mutex.Lock()
	fmt.Println("The lock is locked. (main)")


	//并发的开启3个g
	for i := 1; i <= 3; i++ {
		go func(i int) {
			//gi准备将锁锁起来了，但是main还未解锁额，除非解锁了，才会抢到，否则会阻塞
			//另外就算main中解锁了，每次只会有一个g抢到锁，以此类推
			fmt.Printf("Lock the lock. (g%d)\n", i)
			mutex.Lock()
			fmt.Printf("The lock is locked. (g%d)\n", i)
		}(i)
	}


    //main故意小睡了一会，就是看看g1 g2 g3是否阻塞
	time.Sleep(time.Second)
	//main准备解锁了
	fmt.Println("Unlock the lock. (main)")
	mutex.Unlock()
	fmt.Println("The lock is unlocked. (main)")
	//main又故意睡了一会就是看看g1 g2 g3是否不阻塞了
	time.Sleep(time.Second)
}

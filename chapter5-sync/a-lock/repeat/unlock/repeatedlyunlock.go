package main

import (
	"fmt"
	"sync"
)

func main() {
	defer func() {
		fmt.Println("Try to recover the panic.")
		//从Go1.8 开始此类重复解锁导致的运行时恐慌变为了不可以恢复，所以您是捕获不到的
		if p := recover(); p != nil {
			fmt.Printf("Recovered the panic(%#v).\n", p)
		}
	}()
	var mutex sync.Mutex
	//上锁
	fmt.Println("Lock the lock.")
	mutex.Lock()
	fmt.Println("The lock is locked.")
	fmt.Println("Unlock the lock.")
	//解锁
	mutex.Unlock()
	fmt.Println("The lock is unlocked.")
	//重复解锁
	fmt.Println("Unlock the lock again.")
	mutex.Unlock()
}

package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {

	//接收通道1
	sigRecv1 := make(chan os.Signal, 1)
	sigs1 := []os.Signal{syscall.SIGINT, syscall.SIGQUIT}
	signal.Notify(sigRecv1, sigs1...)//只接收SIGINT   SIGQUIT
	//接收通道2
	sigRecv2 := make(chan os.Signal, 1)
	sigs2 := []os.Signal{syscall.SIGQUIT}
	signal.Notify(sigRecv2, sigs2...)//只接收SIGQUIT


	//=========================并发处理接收到的信号
	var wg sync.WaitGroup
	wg.Add(2)
	//接收通道1的读取
	go func() {
		for sig := range sigRecv1 {
			fmt.Printf("Received a signal from sigRecv1: %s\n", sig)
		}
		//只有该通道被关闭了，才会执行到这里额
		fmt.Printf("End. [sigRecv1]\n")
		wg.Done()
	}()
	//接收通道2的读取
	go func() {
		for sig := range sigRecv2 {
			fmt.Printf("Received a signal from sigRecv2: %s\n", sig)
		}
		fmt.Printf("End. [sigRecv2]\n")
		wg.Done()
	}()
	wg.Wait()





}

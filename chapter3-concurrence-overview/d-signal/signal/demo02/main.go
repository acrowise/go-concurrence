package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

// 监听指定信号
func main()  {
	//合建chan
	c := make(chan os.Signal)
	//监听指定信号 ctrl+c kill
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGUSR1, syscall.SIGUSR2)
	//阻塞直到有信号传入
	fmt.Println("启动")
	//阻塞直至有信号传入
	s := <-c
	fmt.Println("退出信号", s)
}


/*

启动
go run example-2.go
启动

ctrl+c退出,输出
退出信号 interrupt

kill pid 输出
退出信号 terminated

kill -USR1 pid 输出
退出信号 user defined signal 1

kill -USR2 pid 输出
退出信号 user defined signal 2


 */
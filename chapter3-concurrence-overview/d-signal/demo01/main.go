package main

import (
	"os"
	"os/signal"
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



	
}

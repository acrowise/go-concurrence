package main

import (
	"fmt"
	"os/exec"
	"bufio"
)

func main() {


	//第一步:创建命令行
	s:="My first command comes from golang."
	cmd0 := exec.Command("echo", "-n",s)
	//第二步:创建输出管道
	stdout0, err := cmd0.StdoutPipe() //io.ReadCloser
	if err != nil {
		fmt.Printf("Error: Couldn't obtain the stdout pipe for command No.0: %s\n", err)
		return
	}
	//第三步:执行命令
	if err := cmd0.Start(); err != nil {
		fmt.Printf("Error: The command No.0 can not be startup: %s\n", err)
		return
	}
	//第四步:直接使用带缓冲的读取器(这种与我们sdk中读取的方式是一样的，是bufio.Reader中的Read方法，里面自封装了缓冲机制)
	//直接传递输出管道，创建一个缓冲读取器
	outputBuf0 := bufio.NewReader(stdout0)
	//命令的输出本就是一行内容，所以不需要循环判断读取，直接读取一行就over了，所以不需要判断还有没有剩余行额，具体的细节判断封装在ReadLine()内部
	output0, _, err := outputBuf0.ReadLine()
	if err != nil {
		fmt.Printf("Error: Couldn't read data from the pipe: %s\n", err)
		return
	}

	fmt.Printf("命令输出到管道中的内容为: %s\n",string(output0))








}

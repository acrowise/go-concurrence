package main

import (
	"fmt"
	"os/exec"
	"bufio"
)

func main() {


	cmd0 := exec.Command("/bin/bash", "-c", "echo sam is a good man.")

	//为了创建一个能够获取命令行的输出管道，需要在Start之前
	//StdoutPipe方法返回一个在命令Start后与命令标准输出关联的管道。Wait方法获知命令结束后会关闭这个管道，一般不需要显式的关闭该管道。
	//但是在从管道读取完全部数据之前调用Wait是错误的；同样使用StdoutPipe方法时调用Run函数也是错误的
	stdout0, err := cmd0.StdoutPipe()
	if err != nil {
		fmt.Printf("Error: Couldn't obtain the stdout pipe for command No.0: %s\n", err)
		return
	}
	defer stdout0.Close()

	//Start开始执行cmd0包含的命令，但并不会等待该命令完成即返回。Wait方法会返回命令的返回状态码并在命令返回后释放相关的资源。
	if err := cmd0.Start(); err != nil {
		fmt.Printf("Error: The command No.0 can not be startup: %s\n", err)
		return
	}
	//一个更加快捷的方法是，一开始就使用带缓冲的读取器，从输出管道中读取数据
	outputBuf0 := bufio.NewReader(stdout0)
	output0, _, err := outputBuf0.ReadLine()
	if err != nil {
		fmt.Printf("Error: Couldn't read data from the pipe: %s\n", err)
		return
	}
	fmt.Printf("命令输出到管道中的内容为: %s\n", string(output0))



}

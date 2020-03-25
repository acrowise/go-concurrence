package main

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"
)

func main() {


	//cmd0 := exec.Command("/bin/bash", "-c", "echo sam is a good man.")
	cmd0 := exec.Command("/bin/bash", "-c", "sleep 5;echo sam is a good man.")

	//为了创建一个能够获取命令行的输出管道，需要在Start之前创建一个管道，这个管道的作用就是把左边的命令的标准输出捕获起来，以供右边命令的标准输入使用
	//StdoutPipe方法返回一个在命令Start后与命令标准输出关联的管道。
	//Wait方法获知命令结束后会关闭这个管道，一般不需要显式的关闭该管道。
	//但是在从管道读取完全部数据之前调用Wait是错误的；同样使用StdoutPipe方法时调用Run函数也是错误的，即使用StdoutPipe是需要使用Start方法执行命令的
	stdout0, err := cmd0.StdoutPipe() //io.ReadCloser
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

	//读取命令行的输出值
	//stdout0的类型是io.ReadCloser，有了它就可以调用它的Read方法来获取命令的输出了
	//这里的for达到了一种wait的作用，这样就能保障上述执行完毕了
	var outputBuf0 bytes.Buffer
	for {
		//为了达到效果，我故意把字节切片tempOutput的长度设置的很小
		tempOutput := make([]byte, 5)
		n, err := stdout0.Read(tempOutput)
		if err != nil {
			if err == io.EOF {
				break
			} else {
				fmt.Printf("Error: Couldn't read data from the pipe: %s\n", err)
				return
			}
		}
		if n > 0 {
			//为了收集每次迭代读到的输出内容，它们会依次被存放到一个缓冲区outputBuf0中
			outputBuf0.Write(tempOutput[:n])
		}
	}
	fmt.Printf("命令输出到管道中的内容为: %s\n", outputBuf0.String())

}

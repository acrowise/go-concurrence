package main

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"
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
	//第四步:借用缓冲区，循环读取，直到输出管道中没有可以读取的数据了
	//这种就是比较传统的读取方式了,本质就是io.Reader中的Read方法，只是每次读取之后放到了我们引入的缓冲bytes.Buffer中罢了
	var outputBufo bytes.Buffer

	for{
		 tempOutput := make([]byte,5)
		 n,err := stdout0.Read(tempOutput)
		 if err != nil{
		 	 if err == io.EOF {
		 	 	break
			 }else {
			 	 fmt.Printf("Error: Could't read data from the pipe: %s\n",err)
				 return
			 }
		 }
		 if n>0{
		 	 outputBufo.Write(tempOutput[:n])
		 }


	}



	fmt.Printf("命令输出到管道中的内容为: %s\n", string(outputBufo.String()))








}

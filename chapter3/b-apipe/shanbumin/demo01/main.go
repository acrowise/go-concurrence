package main

import (
	"fmt"
	"os/exec"
)

func main() {


	//第一步:创建命令行
	s:="My first command comes from golang."
	cmd0 := exec.Command("echo", "-n",s)
	//第二步:创建输出管道
	//本质这里得到的管道是通过os.Pipe函数生成的。只不过，该方法内部又对生成的管道做了少许的附加处理
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
	//第四步:简单读取
	//@todo  简单读取必须要求承载读取内容的切片的len足够大,大过s的长度才不会读丢失了
	output0:=make([]byte,30)
	//n的值如果输出是30,那么基本上是没有读完，自身已经盛满了，除非恰好读取的s正好也是30的长度
	n,err:=stdout0.Read(output0)
	if err!=nil{
		fmt.Printf("Error:Couldn't read data from the pipe:%s\n",err)
		return
	}

	fmt.Println("err",err) //当stdout0中已经没有数据的时候再次读取的时候会返回io.EOF,否则只要有数据存在，读取成功必然返回nil
	fmt.Println("我们要读取的字节数为",len(s))
	fmt.Println("一次性读取的返回值为：",n)
	fmt.Printf("命令输出到管道中的内容为: %s\n", string(output0))








}

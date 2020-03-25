package main

import (
	"fmt"
	"bytes"
	"os/exec"
)


//这个基于cmd1和cmd2的示例模拟出了操作系统命令     ps aux|grep   doing
func main() {

	//第一步:cmd1  cmd2
	cmd1 := exec.Command("ps", "aux")
	cmd2 := exec.Command("grep", "doing")

	//第二步:设置cmd1的Stdout字段，然后启动cmd1，并等待它运行完毕
	//因为*bytes.Buffer类型实现了io.Writer接口，所以我才能把&outputBuf1赋给cmd1.Stdout。这样命令cmd1启动后的所有输出内容就都会被写入到outputBuf1
	var outputBuf1 bytes.Buffer
	cmd1.Stdout = &outputBuf1 //将cmd1的标准输出与&outputBuf1关联起来
	if err := cmd1.Start(); err != nil {
		fmt.Printf("Error: The first command can not be startup %s\n", err)
		return
	}
	if err := cmd1.Wait(); err != nil {
		fmt.Printf("Error: Couldn't wait for the first command: %s\n", err)
		return
	}

	//第三步:接下来，再设置cmd2的Stdin和Stdout字段，启动cmd2，并等待它运行完毕:
	cmd2.Stdin = &outputBuf1 //将cmd2的标准输入与&outputBuf1关联起来
	var outputBuf2 bytes.Buffer
	cmd2.Stdout = &outputBuf2 //将cmd2的标准输出与&outputBuf2关联起来
	if err := cmd2.Start(); err != nil {
		fmt.Printf("Error: The second command can not be startup: %s\n", err)
		return
	}
	if err := cmd2.Wait(); err != nil {
		fmt.Printf("Error: Couldn't wait for the second command: %s\n", err)
		return
	}
	fmt.Printf("%s\n", outputBuf2.Bytes())
}

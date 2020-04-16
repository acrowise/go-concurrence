package main

import (
	"fmt"
	"os/exec"
)

func main() {



	//************************************** 上述可以通过下面的快速方式实现,这是这个是阻塞的而已，阻塞应该更适合命令行执行

	var (
		cmd0 *exec.Cmd
		output []byte
		err error
	)
	cmd0 = exec.Command("/bin/bash", "-c", "echo sam is a good man.")

	// 执行了命令, 同时捕获了子进程的输出( pipe )
	if output,err = cmd0.CombinedOutput(); err != nil {
		fmt.Printf("Error: The command No.0 can not be startup: %s\n", err)
		return
	}
	// 打印子进程的输出
	fmt.Printf("命令输出到管道中的内容为: %s\n", string(output))


}

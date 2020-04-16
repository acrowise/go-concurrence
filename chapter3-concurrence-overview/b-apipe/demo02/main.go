package main

import (
	"fmt"
	"os/exec"
)

func main() {


	cmd0 := exec.Command("/bin/bash", "-c", "sleep 5")

	//Run执行cmd0包含的命令，并阻塞直到完成。
	//如果命令成功执行，stdin、stdout、stderr的转交没有问题，并且返回状态码为0，方法的返回值为nil；
	//如果命令没有执行或者执行失败，会返回*ExitError类型的错误；否则返回的error可能是表示I/O问题。
	if err := cmd0.Run(); err != nil {
		fmt.Printf("Error: The command No.0 can not be startup: %s\n", err)
		return
	}





	
}

package main

import (
	"fmt"
	"os/exec"
)

func main() {


	cmd0 := exec.Command("/bin/bash", "-c", "sleep 5")

	//Start开始执行cmd0包含的命令，但并不会等待该命令完成即返回。Wait方法会返回命令的返回状态码并在命令返回后释放相关的资源。
	if err := cmd0.Start(); err != nil {
		fmt.Printf("Error: The command No.0 can not be startup: %s\n", err)
		return
	}


	//当使用了Start执行命令之后，Wait方法可以确保它执行完毕，二者组合与Run的效果等价额
	//if err := cmd0.Wait(); err != nil {
	//	fmt.Printf("Error: Couldn't wait for the first command: %s\n", err)
	//	return
	//}





	
}

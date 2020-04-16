package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime/debug"
	"strconv"
	"strings"
	"errors"
	"bytes"
	"syscall"
)

func main() {
	//捕获异常
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("Fatal Error: %s\n", err)
			debug.PrintStack()
		}
	}()
	//初始化管道命令
	// ps aux | grep "show" | grep -v "grep" | grep -v "go run" | awk '{print $2}'  ==>能够获取show.go运行的进程id号
	cmds := []*exec.Cmd{
		exec.Command("ps", "aux"),
		exec.Command("grep", "show"),
		exec.Command("grep", "-v", "grep"),
		exec.Command("grep", "-v", "go run"),
		exec.Command("awk", "{print $2}"),
	}
	output, err := runCmds(cmds)//管道形式执行上述初始化的命令
	if err != nil {
		fmt.Printf("Command Execution Error: %s\n", err)
		return
	}
	//fmt.Println(output)
	//获取执行命令的各个pid值
	pids, err := getPids(output)
	if err != nil {
		fmt.Printf("PID Parsing Error: %s\n", err)
		return
	}
	fmt.Printf("Target PID(s):\n%v\n", pids)
	//遍历给这些进程pids分别发送SIGQUIT信号
	for _, pid := range pids {
		//查找进程，该函数返回一个*os.Process类型的值
		proc, err := os.FindProcess(pid)
		if err != nil {
			fmt.Printf("Process Finding Error: %s\n", err)
			return
		}
		sig := syscall.SIGQUIT
		fmt.Printf("Send signal '%s' to the process (pid=%d)...\n", sig, pid)
		//调用进程值的Signal方法，可以向该进程发送一个信号,这个方法接受一个os.Signal类型的参数值并会返回一个error类型的值
		err = proc.Signal(sig)
		if err != nil {
			fmt.Printf("Signal Sending Error: %s\n", err)
			return
		}
	}
}

//ps aux | grep "show" | grep -v "grep" | grep -v "go run" | awk '{print $2}'
func runCmds(cmds []*exec.Cmd) ([]string, error) {
	//第一步:参数检验
	if cmds == nil || len(cmds) == 0 {
		return nil, errors.New("The cmd slice is invalid!")
	}
	//第二步:循环通过管道 执行命令行
	first := true
	var output []byte
	var err error
	for _, cmd := range cmds {
		fmt.Printf("Run command: %v\n", getCmdPlaintext(cmd))
		//如果不是第一个执行的命令，则需要将上一个命令的标准输出读入到该命令的标准输入中
		if !first {
			var stdinBuf bytes.Buffer
			stdinBuf.Write(output)
			cmd.Stdin = &stdinBuf
		}
		var stdoutBuf bytes.Buffer
		cmd.Stdout = &stdoutBuf
		if err = cmd.Start(); err != nil {
			return nil, getError(err, cmd)
		}
		if err = cmd.Wait(); err != nil {
			return nil, getError(err, cmd)
		}
		output = stdoutBuf.Bytes()
		//fmt.Printf("Output:\n%s\n", string(output))
		if first {
			first = false
		}
	}
	//第三步: 上述循环完毕后，必然是最后一个命令的标准输出内容了(有可能捕获到了多个含有show字样的进程额，所以放到切片lines中)
	var lines []string
	var outputBuf bytes.Buffer
	outputBuf.Write(output)
	for {
		line, err := outputBuf.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return nil, getError(err, nil)
			}
		}
		lines = append(lines, string(line))
	}
	return lines, nil
}


//将切片中的string编程int类型的元素
func getPids(strs []string) ([]int, error) {
	var pids []int
	for _, str := range strs {
		pid, err := strconv.Atoi(strings.TrimSpace(str))
		if err != nil {
			return nil, err
		}
		pids = append(pids, pid)
	}
	return pids, nil
}


//判断某个命令的绝对路径 如 ps aux ===>  /bin/ps aux
func getCmdPlaintext(cmd *exec.Cmd) string {
	var buf bytes.Buffer
	buf.WriteString(cmd.Path)
	for _, arg := range cmd.Args[1:] {
		buf.WriteRune(' ')
		buf.WriteString(arg)
	}
	return buf.String()
}

//错误统一封装报错格式，以及追加部分信息
func getError(err error, cmd *exec.Cmd, extraInfo ...string) error {
	var errMsg string
	if cmd != nil {
		errMsg = fmt.Sprintf("%s  [%s %v]", err, (*cmd).Path, (*cmd).Args)
	} else {
		errMsg = fmt.Sprintf("%s", err)
	}
	if len(extraInfo) > 0 {
		errMsg = fmt.Sprintf("%s (%v)", errMsg, extraInfo)
	}
	return errors.New(errMsg)
}
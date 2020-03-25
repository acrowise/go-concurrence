package main

import (
	"fmt"
	"os"
	"time"
)

//普通的命名管道操作
func main() {

	//这里是os包中的Pipe
	//Go使用系统函数来创建命名管道,并把它的两端封装成两个*os.File类型的值
	//reader代表了该管道输出端的*os.File类型的值,因为管道是单向的，所以不能调用reader.Write
	//writer代表了该管道输入端的*os.File类型的值,因为管道是单向的，所以不能调用writer.Read
	reader, writer, err := os.Pipe()
	if err != nil {
		fmt.Printf("Error: Couldn't create the named pipe: %s\n", err)
	}

	//这里的读和写要并发运行，为什么强调是并发运行呢？
	//因为命名管道默认会在其中一端还未就绪的时候阻塞另一端的进程。Go提供给我们的命名管道的行为特征也是如此。
	//所以，如果顺序执行这两段代码，那么程序肯定会被永远阻塞在语句***处(具体阻塞在哪里，看哪个先调用求值喽)
	//读
	go func() {
		output := make([]byte, 100)
		n, err := reader.Read(output)
		if err != nil {
			fmt.Printf("Error: Couldn't read data from the named pipe: %s\n", err)
		}
		fmt.Printf("Read %d byte(s). [file-based pipe]\n", n)
	}()
	//写
	input := make([]byte, 26)
	for i := 65; i <= 90; i++ {
		input[i-65] = byte(i) //转成byte类型，放入input切片中
	}
	fmt.Println(input)
	n, err := writer.Write(input)
	if err != nil {
		fmt.Printf("Error: Couldn't write data to the named pipe: %s\n", err)
	}
	fmt.Printf("Written %d byte(s). [file-based pipe]\n", n)
	time.Sleep(200 * time.Millisecond)

}

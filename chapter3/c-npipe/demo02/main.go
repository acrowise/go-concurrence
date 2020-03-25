package main

import (
	"fmt"
	"io"
	"time"
)

//一个基于内存的有原子性操作保证的命名管道(多路复用，如存在多个输入端同时写入数据的时候，就需要考虑管道的原子性的问题，因为操作系统提供的管道是不提供原子操作支持的。)
//为此，Go在标准库代码包io中提供了一个基于内存的有原子性操作保证的管道(内存管道)
func main() {

	//这里是io包中的Pipe而不再是os包的Pipe了
	//reader代表了该管道输出端的*io.PipeReader类型的值,该类型对比上面的就更有约束感了，直接就限制只能进行Read，更严格保证了管道的单向性
	//writer代表了该管道输入端的*io.PipeWriter类型的值
	reader, writer := io.Pipe()

	//读
	go func() {
		output := make([]byte, 100)
		n, err := reader.Read(output)
		if err != nil {
			fmt.Printf("Error: Couldn't read data from the named pipe: %s\n", err)
		}
		fmt.Printf("Read %d byte(s). [in-memory pipe]\n", n)
	}()

	//写
	input := make([]byte, 26)
	for i := 65; i <= 90; i++ {
		input[i-65] = byte(i)
	}
	n, err := writer.Write(input)
	if err != nil {
		fmt.Printf("Error: Couldn't write data to the named pipe: %s\n", err)
	}
	fmt.Printf("Written %d byte(s). [in-memory pipe]\n", n)
	time.Sleep(200 * time.Millisecond)

}

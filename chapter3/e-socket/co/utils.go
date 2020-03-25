package co

import (
	"fmt"
	"math"
	"net"
	"strconv"
	"bytes"
	"strings"
)

const (
	DELIMITER      = '\t' //表示一个作为数据边界的单字节字符
)



//**********封装了几个打印输出的方法*****************
func printLog(role string, sn int, format string, args ...interface{}) {
	if !strings.HasSuffix(format, "\n") {
		format += "\n"
	}
	fmt.Printf("%s[%d]: %s", role, sn, fmt.Sprintf(format, args...))
}

func PrintServerLog(format string, args ...interface{}) {
	printLog("Server", 0, format, args...)
}

func PrintClientLog(sn int, format string, args ...interface{}) {
	printLog("Client", sn, format, args...)
}



//****************************************  数据处理辅助函数  ******************************
//尝试将数据转换为int32类型
func StrToInt32(str string) (int32, error) {
	num, err := strconv.ParseInt(str, 10, 0)
	if err != nil {
		return 0, fmt.Errorf("\"%s\" is not integer", str)
	}
	if num > math.MaxInt32 || num < math.MinInt32 {
		return 0, fmt.Errorf("%d is not 32-bit integer", num)
	}
	return int32(num), nil
}
//计算立方根
func Cbrt(param int32) float64 {
	return math.Cbrt(float64(param))
}


//**********************连接间的读与写**************

// 千万不要使用这个版本的read函数！ 这段代码有陷阱
// 作为一个处在TCP/IP协议栈的应用层的程序，负责切分数据并生成有实际意义的消息。
// 即使在最简单的情况下，应用程序也知道怎样在接收到的字节流上进行切分，你可以按照自己的要求去编写实现切分操作的程序
// 不过，还有一个更简单的方法：利用标准库代码包bufio中的API实现一些较复杂的数据切分操作。
//func read(conn net.Conn) (string, error) {
//	reader := bufio.NewReader(conn)
//	readBytes, err := reader.ReadBytes(DELIMITER)
//	if err != nil {
//		return "", err
//	}
//	return string(readBytes[:len(readBytes)-1]), nil
//}


//故意每次都只读一个字节
//当前底层没有帮我们做到屏蔽部分读的特性，将数据一次性读取完毕返回给我们
//这里需要我们自己判断什么时候读结束，这是合理的，比如遇到DELIMITER 我们约定好的结束符，我们就认为读取该消息结束了
func Read(conn net.Conn) (string, error) {
	readBytes := make([]byte, 1)
	var buffer bytes.Buffer
	for {
		_, err := conn.Read(readBytes)
		if err != nil {
			return "", err
		}
		readByte := readBytes[0]
		if readByte == DELIMITER {
			break
		}
		buffer.WriteByte(readByte)
	}
	return buffer.String(), nil
}

//这几行代码较完整地展现了一个在TCP连接之上读取数据的流程
func Read2(conn net.Conn) (string, error) {
	//声明了一个bytes.Buffer类型值，并以此来存储接收到的所有数据
	var dataBuffer bytes.Buffer
	b := make([]byte, 10)
	//无限循环读取
	//总是先在变量conn的值上调用Read方法以读取从网络上接收到的数据，并在确定未发生任何错误之后把数据追加到dataBuffer的值中
	for {
		n, err := conn.Read(b)
		if err != nil {
			return "", err
		}
		readByte := b[0]
		if readByte == DELIMITER {
			break
		}
		dataBuffer.Write(b[:n])
	}
	return dataBuffer.String(), nil
}




//写数据
//go socket编程API程序在底层帮我们屏蔽了非阻塞式socket接口的部分写特性
//相关API直到把所有数据全部写入到socket的发送缓冲区之后才会返回，除非写的过程中发生了错误。
//所以对于写入，我们不需要在上层做任务判断处理
func Write(conn net.Conn, content string) (int, error) {
	var buffer bytes.Buffer
	buffer.WriteString(content)
	buffer.WriteByte(DELIMITER)
	return conn.Write(buffer.Bytes())
}
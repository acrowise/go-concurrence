package main

import (
	"fmt"
	"gopcp.v2/chapter3/e-socket/co"
	"io"
	"net"
	"time"
)

const (
	SERVER_NETWORK = "tcp"
	SERVER_ADDRESS = "127.0.0.1:8085"
	DELIMITER      = '\t' //表示一个作为数据边界的单字节字符
)

//接收客户端的请求，计算请求数据的立方根，并把对结果的描述返回给客户端程序。
//这种按流传递的方式，客户端和服务器端一定要约好数据边界的格式，好利于切分啊，比如 \t符号
//1.net.Listener()
//2.listener.Accept()

func main() {

	//①根据给定的网络协议和地址创建一个监听器
	listener, err := net.Listen(SERVER_NETWORK, SERVER_ADDRESS)
	if err != nil {
		fmt.Printf("Listen Error: %s\n", err)
		return
	}
	defer listener.Close()
	fmt.Println("Got listener for the server. (local address: %s)", listener.Addr())
	//②一旦成功获得监听器，就可以开始等待客户端的连接请求了
	for {
		conn, err := listener.Accept()// 阻塞直至新连接到来。
		if err != nil {
			fmt.Println("Accept Error: %s", err)
			continue
		}
		// 返回远端网络地址
		fmt.Println("Established a connection with a client application. (remote address: %s)", conn.RemoteAddr())
		//采用并发的方式处理连接
		go handleConn(conn)
	}
}

//具体处理客户端的连接请求
func handleConn(conn net.Conn) {
	defer func() {
		conn.Close()
	}()

	for {
		//设定该连接的读操作deadline，参数为零值表示不设置期限
		//设置这个能起到关闭闲置连接功能，超时错误的发生意味着当前连接大部分可以判定为闲置连接了
		conn.SetReadDeadline(time.Now().Add(10 * time.Second))
		//读取，从连接中读取一段以数据分界符为结尾的数据
		strReq, err := co.Read(conn)
		if err != nil {
			if err == io.EOF {
				fmt.Println("The connection is closed by another side.")
			} else {
				fmt.Println("Read Error: %s", err)
			}
			break
		}
		fmt.Printf("Received request: %s.\n", strReq)
		//这部分代码实现的功能是检查数据块是否可以转换为一个int32类型的值，如果能，就立即计算它的立方根，否则就向客户端程序发送一条错误信息
		intReq, err := co.StrToInt32(strReq)
		if err != nil {
			n, err := co.Write(conn, err.Error())
			fmt.Printf("Sent error message (written %d bytes): %s.\n", n, err)
			continue
		}
		floatResp := co.Cbrt(intReq)
		//写数据
		respMsg := fmt.Sprintf("The cube root of %d is %f.\n", intReq, floatResp)
		n, err := co.Write(conn, respMsg)
		if err != nil {
			fmt.Printf("Write Error: %s\n", err)
		}
		fmt.Printf("Sent response (written %d bytes): %s.\n", n, respMsg)
	}
}




















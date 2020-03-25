package main

import (
	"fmt"
	"gopcp.v2/chapter3/e-socket/co"
	"io"
	"math/rand"
	"net"
	"time"
)

const (
	SERVER_NETWORK = "tcp"
	SERVER_ADDRESS = "127.0.0.1:8085"
	DELIMITER      = '\t' //表示一个作为数据边界的单字节字符
)




func clientGo(id int) {

	//①试图与服务端程序建立连接
	//Dial方法固定了连接超时时间(具体取决于操作系统，比如在Linux操作系统内核中，把基于TCP协议的连接请求的超时时间设定为75秒)
	//conn, err := net.Dial(SERVER_NETWORK, SERVER_ADDRESS)
	//第一个参数,由于是客户端，我们不像服务器端那样只能填写tcp之类的流协议，这里可选值更宽广一些，因为在发送数据之前不一定要先建立连接
	//所以可以是udp  ip ...这些面向无连接型的协议都是可以的
	//即该方法可以与tcp或者upd服务器通讯，目前我们的服务端开启的是tcp协议，所以这里我们填写tcp
	conn, err := net.DialTimeout(SERVER_NETWORK, SERVER_ADDRESS, 2*time.Second)
	if err != nil {
		co.PrintClientLog(id, "Dial Error: %s", err)
		return
	}
	defer conn.Close()
	//返回本地网络地址  返回远端网络地址
	//LocalAddr指的是当前程序所使用的地址，即客户端自己的地址(客户端使用的端口号可以由应用程序指定，也可以由系统内核动态分配，目前我们使用后者)
	//RemoteAddr()是参与通信的另一端所使用的地址,即服务器端地址
	co.PrintClientLog(id, "Connected to server. (remote address: %s, local address: %s)", conn.RemoteAddr(), conn.LocalAddr())
	time.Sleep(200 * time.Millisecond)

	//②发送请求数据
	//确定把每个客户端发送的请求数据块的数量定为5个。
	requestNumber := 5
	//设定当前连接上的I/O操作的超时时间
	conn.SetDeadline(time.Now().Add(5 * time.Millisecond))
	for i := 0; i < requestNumber; i++ {
		req := rand.Int31() //可以随机生成一个int32类型值
		n, err := co.Write(conn, fmt.Sprintf("%d", req))
		if err != nil {
			co.PrintClientLog(id, "Write Error: %s", err)
			continue
		}
		co.PrintClientLog(id, "Sent request (written %d bytes): %d.", n, req)
	}
	//③准备接收响应数据块
	for j := 0; j < requestNumber; j++ {
		strResp, err := co.Read(conn)
		if err != nil {
			if err == io.EOF {
				co.PrintClientLog(id, "The connection is closed by another side.")
			} else {
				co.PrintClientLog(id, "Read Error: %s", err)
			}
			break
		}
		co.PrintClientLog(id, "Received response: %s.", strResp)
	}
}

func main() {

	go clientGo(1)

	time.Sleep(10 * time.Second)
}

package testhelper

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"time"

	loadgenlib "shanbumin/go-concurrence/chapter4-goroutine/k-loadgen/lib"
)

const (
	DELIM = '\n' // 分隔符。
)

// operators 代表操作符切片。
var operators = []string{"+", "-", "*", "/"}

// TCPComm 表示TCP通讯器的结构。
//TCPComm就是我们此次要实现的Caller接口类型
type TCPComm struct {
	addr string
}

// NewTCPComm 会新建一个TCP通讯器。
func NewTCPComm(addr string) loadgenlib.Caller {
	return &TCPComm{addr: addr}
}

// BuildReq 会构建一个请求。
func (comm *TCPComm) BuildReq() loadgenlib.RawReq {
	id := time.Now().UnixNano()
	//请求发送的数据的构建
	sreq := ServerReq{
		ID: id,
		Operands: []int{
			int(rand.Int31n(1000) + 1),
			int(rand.Int31n(1000) + 1)},
		Operator: func() string {
			return operators[rand.Int31n(100)%4]
		}(),
	}
	bytes, err := json.Marshal(sreq)//结构体===>json文本的[]类型值
	if err != nil {
		panic(err)
	}
	rawReq := loadgenlib.RawReq{ID: id, Req: bytes}
	return rawReq
}

// Call 会发起一次通讯。
//第二个参数是超时判断，不过由于载荷器内部已经做了统一的判断所以不是强制的。
func (comm *TCPComm) Call(req []byte, timeoutNS time.Duration) ([]byte, error) {
	//建立tcp网络通信连接
	conn, err := net.DialTimeout("tcp", comm.addr, timeoutNS)
	if err != nil {
		return nil, err
	}
	//如果通信建立成功，就先将请求数据写入连接
	_, err = write(conn, req, DELIM)
	if err != nil {
		return nil, err
	}
	//然后在成功后再等待并试着从连接中读取影响
	return read(conn, DELIM)
}

// CheckResp 会检查响应内容。
// 调用该方法的地方是载荷发生器接收到被测软件的响应之后。
// 如果原始响应中没有携带任何错误,那么载荷发生器就会调用它对原始响应进行进一步的检查,并根据检查结果设置其返回的调用结果的Code字段值和Msg字段值。
func (comm *TCPComm) CheckResp(rawReq loadgenlib.RawReq, rawResp loadgenlib.RawResp) *loadgenlib.CallResult {
	//在开始处,需要先对调用结果进行必要的初始化
	var commResult loadgenlib.CallResult
	commResult.ID = rawResp.ID
	commResult.Req = rawReq
	commResult.Resp = rawResp
	//请求数据转对应的结构体(方便检查)
	var sreq ServerReq
	err := json.Unmarshal(rawReq.Req, &sreq)
	if err != nil {
		commResult.Code = loadgenlib.RET_CODE_FATAL_CALL
		commResult.Msg = fmt.Sprintf("Incorrectly formatted Req: %s!\n", string(rawReq.Req))
		return &commResult
	}
	//响应数据转对应的结构体(方便检查)
	var sresp ServerResp
	err = json.Unmarshal(rawResp.Resp, &sresp)
	if err != nil {
		commResult.Code = loadgenlib.RET_CODE_ERROR_RESPONSE
		commResult.Msg = fmt.Sprintf("Incorrectly formatted Resp: %s!\n", string(rawResp.Resp))
		return &commResult
	}
    //对sresp正确性的检查
	//(1)原始响应是否与该原始请求相对应;如果不是，该项检查未通过
	if sresp.ID != sreq.ID {
		commResult.Code = loadgenlib.RET_CODE_ERROR_RESPONSE
		commResult.Msg = fmt.Sprintf("Inconsistent raw id! (%d != %d)\n", rawReq.ID, rawResp.ID)
		return &commResult
	}
	//(2)被测软件在处理请求的过程中是否发生了错误
	if sresp.Err != nil {
		commResult.Code = loadgenlib.RET_CODE_ERROR_CALEE
		commResult.Msg = fmt.Sprintf("Abnormal server: %s!\n", sresp.Err)
		return &commResult
	}
	//(3)检查变量sresp的Result字段值是否正确,即它是否为正确的运算结果
	if sresp.Result != op(sreq.Operands, sreq.Operator) {
		commResult.Code = loadgenlib.RET_CODE_ERROR_RESPONSE
		commResult.Msg = fmt.Sprintf("Incorrect result: %s!\n", genFormula(sreq.Operands, sreq.Operator, sresp.Result, false))
		return &commResult
	}


	commResult.Code = loadgenlib.RET_CODE_SUCCESS
	commResult.Msg = fmt.Sprintf("Success. (%s)", sresp.Formula)
	return &commResult
}





//todo 已知,基于TCP协议的通信是使用字节流来传递上层给予的消息的。它会根据具体情况为消息分段,但却无法感知消息的分界。
//todo 因此，需要显式地为请求数据添加结束符,而传给write方法和read方法的参数DELIM就代表了这个结束符,
//todo 这两个方法会用它来切分出单个的请求或响应。

// read 会从连接中读数据直到遇到参数delim代表的字节。
func read(conn net.Conn, delim byte) ([]byte, error) {
	readBytes := make([]byte, 1)
	var buffer bytes.Buffer
	for {
		_, err := conn.Read(readBytes)
		if err != nil {
			return nil, err
		}
		readByte := readBytes[0]
		if readByte == delim {
			break
		}
		buffer.WriteByte(readByte)
	}
	return buffer.Bytes(), nil
}

// write 会向连接写数据，并在最后追加参数delim代表的字节。
func write(conn net.Conn, content []byte, delim byte) (int, error) {
	writer := bufio.NewWriter(conn)
	n, err := writer.Write(content)
	if err == nil {
		writer.WriteByte(delim)
	}
	if err == nil {
		err = writer.Flush()
	}
	return n, err
}

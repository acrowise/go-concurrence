package lib

import (
	"time"
)

//======================================================
//对于同一个载荷而言，其请求、响应和调用结果中，ID字段的值都是一致的，这对于我们了解载荷处理的全过程很有帮助。
// 自定义类型  RawReq 表示原生请求的结构。
type RawReq struct {
	ID  int64
	Req []byte //数据的最底层表现形式就是若干字节，因此将[]byte作为Req字段的类型
}

//自定义类型  RawResp 表示原生响应的结构。
type RawResp struct {
	ID     int64
	Resp   []byte
	Err    error
	Elapse time.Duration
}

//======================================================

// 为了更好地定义调用结果值中的响应代码，定义如下常量:
// RetCode 表示结果代码的类型。
type RetCode int

// 保留 1 ~ 1000 给载荷承受方使用。
const (
	RET_CODE_SUCCESS              RetCode = 0    // 成功。
	RET_CODE_WARNING_CALL_TIMEOUT         = 1001 // 调用超时警告。
	RET_CODE_ERROR_CALL                   = 2001 // 调用错误。
	RET_CODE_ERROR_RESPONSE               = 2002 // 响应内容错误。
	RET_CODE_ERROR_CALEE                  = 2003 // 被调用方（被测软件）的内部错误。
	RET_CODE_FATAL_CALL                   = 3001 // 调用过程中发生了致命错误！
)

// GetRetCodePlain 会依据结果代码返回相应的文字解释。
func GetRetCodePlain(code RetCode) string {
	var codePlain string
	switch code {
	case RET_CODE_SUCCESS:
		codePlain = "Success"
	case RET_CODE_WARNING_CALL_TIMEOUT:
		codePlain = "Call Timeout Warning"
	case RET_CODE_ERROR_CALL:
		codePlain = "Call Error"
	case RET_CODE_ERROR_RESPONSE:
		codePlain = "Response Error"
	case RET_CODE_ERROR_CALEE:
		codePlain = "Callee Error"
	case RET_CODE_FATAL_CALL:
		codePlain = "Call Fatal Error"
	default:
		codePlain = "Unknown result code"
	}
	return codePlain
}

// CallResult 表示调用结果的结构。
//每次的请求结果都是该结构体表示的额
type CallResult struct {
	ID     int64         // ID。 作用是标识调用结果
	Req    RawReq        // 原生请求。
	Resp   RawResp       // 原生响应。
	Code   RetCode       // 响应代码。
	Msg    string        // 结果成因的简述。
	Elapse time.Duration // 耗时。
}





// 声明代表载荷发生器状态的常量。
const (

	STATUS_ORIGINAL uint32 = 0 //代表原始。
	STATUS_STARTING uint32 = 1 //代表正在启动。
	STATUS_STARTED uint32 = 2  //代表已启动。
	STATUS_STOPPING uint32 = 3 //代表正在停止。
	STATUS_STOPPED uint32 = 4  //代表已停止。
)

// Generator 表示载荷发生器的接口。
type Generator interface {

	Start() bool  // 启动载荷发生器。 结果值代表是否已成功启动。
	Stop() bool  // 停止载荷发生器。 结果值代表是否已成功停止。
	Status() uint32  // 获取状态。
	CallCount() int64  // 获取调用计数。每次启动会重置该计数。
}

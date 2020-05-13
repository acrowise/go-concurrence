package loadgen

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"math"
	"sync/atomic"
	"time"
	"shanbumin/go-concurrence/chapter4-goroutine/k-loadgen/lib"
	"shanbumin/go-concurrence/helper/log"
)






// 日志记录器。
var logger = log.DLogger()

// myGenerator 代表载荷发生器的实现类型。
type myGenerator struct {
	//①标准的三个输入参数以及一个结果容器
	timeoutNS   time.Duration        // 处理超时时间，单位：纳秒。
	lps         uint32               // 每秒载荷量。  Loads Per Second
	durationNS  time.Duration        // 负载持续时间，单位：纳秒。
	resultCh    chan *lib.CallResult // 调用结果通道。
	//②加入4个起控制作用的字段
	concurrency uint32               // 载荷并发量。可以认为就是开G的个数啦
	tickets     lib.GoTickets        // Goroutine票池。
	ctx         context.Context      // 上下文。
	cancelFunc  context.CancelFunc   // 取消函数。
    //③载荷发生器有不止一种的状态。状态字段是数值类型的，并且足够短小，还可以用并发安全的方式操作。Go标准库中提供的原子操作方法支持的最短数值类型为int32何uint32
	status      uint32               // 状态。
    //④ 扩展的接口类型
	caller      lib.Caller           // 调用器。
	//⑤ 用于记录调用计数 todo 该类型为什么是包级私有的?难道我们不想让loadgen子包之外的程序使用它吗
	callCount   int64                // 调用计数。


}

// NewGenerator 会新建一个载荷发生器。
// 可以直接传递若干个参数或者一个整合的结构体ParamSet
// 返回值通过接口类型还能起到规范作用,比如返回的gen对应的*myGenerator类型必须实现lib.Generator的方法，否则编译不通过的额
func NewGenerator(pset ParamSet) (lib.Generator, error) {

	logger.Infoln("New a load generator...")
	if err := pset.Check(); err != nil {
		return nil, err
	}
	gen := &myGenerator{
		caller:     pset.Caller,
		timeoutNS:  pset.TimeoutNS,
		lps:        pset.LPS,
		durationNS: pset.DurationNS,
		status:     lib.STATUS_ORIGINAL,
		resultCh:   pset.ResultCh,
	}
	if err := gen.init(); err != nil {
		return nil, err
	}
	return gen, nil
}

// 初始化载荷发生器。
//该方法初始化了另外两个载荷发生器启动前必需的字段——concurrency和tickets
func (gen *myGenerator) init() error {
	var buf bytes.Buffer
	buf.WriteString("Initializing the load generator...")


	//-----------------------  载荷的并发量 ≈ 载荷的响应超时时间 / 载荷的发送间隔时间
	var total64 = int64(gen.timeoutNS)/int64(1e9/gen.lps) + 1
	if total64 > math.MaxInt32 {
		total64 = math.MaxInt32
	}
	gen.concurrency = uint32(total64)
	//goroutine票池的初始化
	tickets, err := lib.NewGoTickets(gen.concurrency)
	if err != nil {
		return err
	}
	gen.tickets = tickets
	//-----------------------

	buf.WriteString(fmt.Sprintf("Done. (concurrency=%d)", gen.concurrency))
	logger.Infoln(buf.String())
	return nil
}

// callOne 会向载荷承受方发起一次调用。
func (gen *myGenerator) callOne(rawReq *lib.RawReq) *lib.RawResp {
	//原子地递增载荷发生器的callCount字段值
	atomic.AddInt64(&gen.callCount, 1)
	//检查参数
	if rawReq == nil {
		return &lib.RawResp{ID: -1, Err: errors.New("Invalid raw request.")}
	}
	//执行调用,记录调用时长
	start := time.Now().UnixNano()
	resp, err := gen.caller.Call(rawReq.Req, gen.timeoutNS)
	end := time.Now().UnixNano()
	elapsedTime := time.Duration(end - start)
	//组装并返回原始响应值
	var rawResp lib.RawResp
	if err != nil {
		errMsg := fmt.Sprintf("Sync Call Error: %s.", err)
		rawResp = lib.RawResp{
			ID:     rawReq.ID,
			Err:    errors.New(errMsg),
			Elapse: elapsedTime}
	} else {
		rawResp = lib.RawResp{
			ID:     rawReq.ID,
			Resp:   resp,
			Elapse: elapsedTime}
	}
	return &rawResp
}




// 一个调用过程分为5个操作步骤:①生成载荷  ②发送载荷并接收响应  ③检查载荷响应  ④生成调用结果   ⑤发送调用结果  (前3个操作步骤都会由使用方在初始化载荷发生器时传入的那个调用器来完成)
//todo  因为对asyncCall方法的每一次调用都意味着会有一个专用goroutine被启用。这里的启用数量由票池控制，我们只需要在适当的时候对g票池中的票进行"获得"和归还操作。
// asyncSend 会异步地调用承受方接口。
func (gen *myGenerator) asyncCall() {
	//取票,当票池中无票可拿时,asyncCall方法所在的goroutine会被阻塞于此。
	gen.tickets.Take()
	go func() {
		defer func() {
			if p := recover(); p != nil {
				//先使用类型断言表达式判断变量p的实际类型是否为error
				err, ok := interface{}(p).(error)
				var errMsg string
				if ok {
					errMsg = fmt.Sprintf("Async Call Panic! (error: %s)", err)
				} else {
					errMsg = fmt.Sprintf("Async Call Panic! (clue: %#v)", p)
				}
				logger.Errorln(errMsg)
				result := &lib.CallResult{
					ID:   -1,
					Code: lib.RET_CODE_FATAL_CALL,
					Msg:  errMsg}
				gen.sendResult(result)
			}
			//还票
			gen.tickets.Return()
		}()

		//步骤① 生成载荷
		rawReq := gen.caller.BuildReq()
		//步骤② 发送载荷并接收响应  (通过callStatus的原子操作巧妙实现超时判断真的是很巧妙呀)
		var callStatus uint32 // 调用状态：0-未调用或调用中；1-调用完成；2-调用超时。
		timer := time.AfterFunc(gen.timeoutNS, func() {
			if !atomic.CompareAndSwapUint32(&callStatus, 0, 2) {
				return
			}
			result := &lib.CallResult{
				ID:     rawReq.ID,
				Req:    rawReq,
				Code:   lib.RET_CODE_WARNING_CALL_TIMEOUT,
				Msg:    fmt.Sprintf("Timeout! (expected: < %v)", gen.timeoutNS),
				Elapse: gen.timeoutNS,
			}
			gen.sendResult(result)
		})
		rawResp := gen.callOne(&rawReq)
		if !atomic.CompareAndSwapUint32(&callStatus, 0, 1) {
			return
		}
		timer.Stop()
		//步骤③  响应处理:就是对原始响应进行再包装，然后再把包装后的响应发给调用结果通道。

		var result *lib.CallResult
		if rawResp.Err != nil {
			result = &lib.CallResult{
				ID:     rawResp.ID,
				Req:    rawReq,
				Code:   lib.RET_CODE_ERROR_CALL,
				Msg:    rawResp.Err.Error(),
				Elapse: rawResp.Elapse}
		} else {
			result = gen.caller.CheckResp(rawReq, *rawResp)
			result.Elapse = rawResp.Elapse
		}
		gen.sendResult(result)
	}()
}






// sendResult 用于发送调用结果。
// 该方法会向调用结果通道发送一个调用结果值。
func (gen *myGenerator) sendResult(result *lib.CallResult) bool {

	//先检查载荷发生器的状态,如果它的状态不是已启动，就不能执行发送操作了。
	if atomic.LoadUint32(&gen.status) != lib.STATUS_STARTED {
		gen.printIgnoredResult(result, "stopped load generator")
		return false
	}
	//若调用结果通道已满,也不能执行发送操作。由于该通道是载荷发生器的使用方传入的,因此无法保证没有这种情况发生。
	//因此,这里需要把发送操作作为一条select语句中的一个case,并添加default分支以确保不会发生阻塞
	select {
	case gen.resultCh <- result:
		return true
	default:
		gen.printIgnoredResult(result, "full result channel")
		return false
	}
}

// printIgnoredResult 打印被忽略的结果。
// 记录未发送的结果
func (gen *myGenerator) printIgnoredResult(result *lib.CallResult, cause string) {
	resultMsg := fmt.Sprintf("ID=%d, Code=%d, Msg=%s, Elapse=%v", result.ID, result.Code, result.Msg, result.Elapse)
	logger.Warnf("Ignored result: %s. (cause: %s)\n", resultMsg, cause)
}

// prepareStop 用于为停止载荷发生器做准备。
//todo 这里为什么会变更两次状态而不一次变更到最终结果呢?
func (gen *myGenerator) prepareToStop(ctxError error) {
	logger.Infof("Prepare to stop load generator (cause: %s)...", ctxError)

	//CAS操作(比较并交换)
	//该方法会先仅在载荷发生器的状态为已启用时,把它变为正在停止状态。
	atomic.CompareAndSwapUint32(&gen.status, lib.STATUS_STARTED, lib.STATUS_STOPPING)

	logger.Infof("Closing result channel...")
	//关闭调用结果通道
	close(gen.resultCh)
	//最后再把状态变为已停止
	atomic.StoreUint32(&gen.status, lib.STATUS_STOPPED)
}

// genLoad 会产生载荷并向承受方发送。
// 该方法总体上控制调用流程的执行
func (gen *myGenerator) genLoad(throttle <-chan time.Time) {

	//使用了一个for循环周期性地向被测软件发送载荷,这个周期的长短由节流阀控制。
	for {
		//------------------ 循环开始处的select ----------------
		select {
		case <-gen.ctx.Done():
			gen.prepareToStop(gen.ctx.Err())
			return
		default:
		}
		//---------------------------------------------------

		//##################
		gen.asyncCall()
		//##################


        //------------------ 循环结尾处的select --------------------
		//如果lps字段的值大于0,就表示节流阀是有效并需要使用的。
		if gen.lps > 0 {
			select {
			case <-throttle:
			case <-gen.ctx.Done():
				gen.prepareToStop(gen.ctx.Err())
				return
			}
		}
		//------------------------------------------------------------


	}
}

//====================== 让*myGenerator类型实现接口lib.Generator
// Start 会启动载荷发生器。
func (gen *myGenerator) Start() bool {
	logger.Infoln("Starting load generator...")
	// 检查是否具备可启动的状态，顺便设置状态为正在启动
	if !atomic.CompareAndSwapUint32(
		&gen.status, lib.STATUS_ORIGINAL, lib.STATUS_STARTING) {
		if !atomic.CompareAndSwapUint32(
			&gen.status, lib.STATUS_STOPPED, lib.STATUS_STARTING) {
			return false
		}
	}

	// 设定节流阀。
	var throttle <-chan time.Time
	if gen.lps > 0 {
		interval := time.Duration(1e9 / gen.lps)
		logger.Infof("Setting throttle (%v)...", interval)
		throttle = time.Tick(interval) //调用Tick函数会返回一个时间类型的接收channel(该结果值是一个可以周期性地传达到期通知的缓冲通道)
	}

	// 初始化上下文和取消函数。
	gen.ctx, gen.cancelFunc = context.WithTimeout(context.Background(), gen.durationNS)

	// 初始化调用计数。
	gen.callCount = 0

	// 设置状态为已启动。
	atomic.StoreUint32(&gen.status, lib.STATUS_STARTED)

	//由于这里启用另一个goroutine来执行生成并发送载荷的流程,因此Start方法是非阻塞的
	go func() {
		// 生成并发送载荷。
		logger.Infoln("Generating loads...")
		gen.genLoad(throttle)
		logger.Infof("Stopped. (call count: %d)", gen.callCount)
	}()
	return true
}

//载荷发生器的手动停止
func (gen *myGenerator) Stop() bool {
	//首先需要检查载荷发生器的状态，若状态不对就直接返回false
	if !atomic.CompareAndSwapUint32(&gen.status, lib.STATUS_STARTED, lib.STATUS_STOPPING) {
		return false
	}
	//执行cancelFunc()字段所代表的方法可以让ctx字段发出停止"信号"
	gen.cancelFunc()
	//需要不断检查状态的变更。如果状态变为已停止，就说明prepareToStop方法已执行完毕，这时就可以返回了。
	for {
		if atomic.LoadUint32(&gen.status) == lib.STATUS_STOPPED {
			break
		}
		time.Sleep(time.Microsecond)
	}
	return true
}

func (gen *myGenerator) Status() uint32 {
	return atomic.LoadUint32(&gen.status)
}

func (gen *myGenerator) CallCount() int64 {
	return atomic.LoadInt64(&gen.callCount)
}

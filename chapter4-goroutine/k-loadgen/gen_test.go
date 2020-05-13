package loadgen

import (
	"testing"
	"time"
	loadgenlib "shanbumin/go-concurrence/chapter4-goroutine/k-loadgen/lib"
	helper "shanbumin/go-concurrence/chapter4-goroutine/k-loadgen/testhelper"
)



// printDetail 代表是否打印详细结果。
var printDetail = false

//todo 用来进行功能测试的函数的名称以"Test"为前缀,并接受一个*testing.T类型的参数

//go  test  -v  -run=TestStart
//go  test  -v  -run=TestStop



//主要对载荷发生器的启动流程、控制和调用流程以及自动停止流程进行了测试
func TestStart(t *testing.T) {

	// 初始化服务器。
	server := helper.NewTCPServer()
	defer server.Close()
	serverAddr := "127.0.0.1:8080"
	t.Logf("Startup TCP server(%s)...\n", serverAddr)
	err := server.Listen(serverAddr)
	if err != nil {
		t.Fatalf("TCP Server startup failing! (addr=%s)!\n", serverAddr)
		t.FailNow()
	}
	// 初始化载荷发生器。
	//传递参数的初始化
	pset := ParamSet{
		Caller:     helper.NewTCPComm(serverAddr),//创建和初始化一个基于TCP协议的调用器
		TimeoutNS:  50 * time.Millisecond, //设定的响应超时时间为50ms
		LPS:        uint32(1000), //每秒载荷量为1000
		DurationNS: 10 * time.Second, //负载持续时间为10s
		ResultCh:   make(chan *loadgenlib.CallResult, 50), //调用结果通道的容量为50
	}
	t.Logf("Initialize load generator (timeoutNS=%v, lps=%d, durationNS=%v)...", pset.TimeoutNS, pset.LPS, pset.DurationNS)
	gen, err := NewGenerator(pset)
	if err != nil {
		t.Fatalf("Load generator initialization failing: %s\n", err)
		t.FailNow()
	}

	// 开始！
	t.Log("Start load generator...")
	gen.Start()  //正常情况是很快就会结束

	// 显示结果。
	countMap := make(map[loadgenlib.RetCode]int)
	//for语句不断尝试从调用结果通道接收结果,并按照响应代码分类对响应计数
	for r := range pset.ResultCh {
		countMap[r.Code] = countMap[r.Code] + 1 //计数
		if printDetail {
			t.Logf("Result: ID=%d, Code=%d, Msg=%s, Elapse=%v.\n", r.ID, r.Code, r.Msg, r.Elapse)
		}
	}

	//展示对调用结果的统计
	var total int //total是所有分类的计数的总和
	t.Log("RetCode Count:")
	for k, v := range countMap {
		codePlain := loadgenlib.GetRetCodePlain(k)
		t.Logf("  Code plain: %s (%d), Count: %d.\n", codePlain, k, v)
		total += v
	}

	t.Logf("Total: %d.\n", total)
	successCount := countMap[loadgenlib.RET_CODE_SUCCESS]
	//被测软件平均每秒有效的处理(或称响应)载荷的数量
	tps := float64(successCount) / float64(pset.DurationNS/1e9)
	t.Logf("Loads per second: %d; Treatments per second: %f.\n", pset.LPS, tps)
}

/* todo  结果通道里存的是这些东西:
type CallResult struct {
	ID     int64         // ID。 作用是标识调用结果
	Req    RawReq        // 原生请求。
	Resp   RawResp       // 原生响应。
	Code   RetCode       // 响应代码。
	Msg    string        // 结果成因的简述。
	Elapse time.Duration // 耗时。
}
*/



//手动停止的测试,比上面就多了该环节而已
func TestStop(t *testing.T) {

	// 初始化服务器。
	server := helper.NewTCPServer()
	defer server.Close()
	serverAddr := "127.0.0.1:8081"
	t.Logf("Startup TCP server(%s)...\n", serverAddr)
	err := server.Listen(serverAddr)
	if err != nil {
		t.Fatalf("TCP Server startup failing! (addr=%s)!\n", serverAddr)
		t.FailNow()
	}

	// 初始化载荷发生器。
	pset := ParamSet{
		Caller:     helper.NewTCPComm(serverAddr),
		TimeoutNS:  50 * time.Millisecond,
		LPS:        uint32(1000),
		DurationNS: 10 * time.Second,
		ResultCh:   make(chan *loadgenlib.CallResult, 50),
	}
	t.Logf("Initialize load generator (timeoutNS=%v, lps=%d, durationNS=%v)...", pset.TimeoutNS, pset.LPS, pset.DurationNS)
	gen, err := NewGenerator(pset)
	if err != nil {
		t.Fatalf("Load generator initialization failing: %s.\n",
			err)
		t.FailNow()
	}

	// 开始！
	t.Log("Start load generator...")
	gen.Start()
	//todo 这里手动停止定时为2s,少于预设的负载持续时间10s。
	timeoutNS := 2 * time.Second
	time.AfterFunc(timeoutNS, func() {
		gen.Stop()
	})

	// 显示调用结果。
	countMap := make(map[loadgenlib.RetCode]int)
	count := 0
	for r := range pset.ResultCh {
		countMap[r.Code] = countMap[r.Code] + 1
		if printDetail {
			t.Logf("Result: ID=%d, Code=%d, Msg=%s, Elapse=%v.\n", r.ID, r.Code, r.Msg, r.Elapse)
		}
		count++
	}

	var total int
	t.Log("RetCode Count:")
	for k, v := range countMap {
		codePlain := loadgenlib.GetRetCodePlain(k)
		t.Logf("  Code plain: %s (%d), Count: %d.\n",
			codePlain, k, v)
		total += v
	}

	t.Logf("Total: %d.\n", total)
	successCount := countMap[loadgenlib.RET_CODE_SUCCESS]
	tps := float64(successCount) / float64(timeoutNS/1e9)
	t.Logf("Loads per second: %d; Treatments per second: %f.\n", pset.LPS, tps)
}

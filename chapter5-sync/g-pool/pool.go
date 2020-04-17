package main

import (
	"fmt"
	"runtime"
	"runtime/debug"
	"sync"
	"sync/atomic"
)
//临时对象池的第二个突出特性是对垃圾回收友好。
//垃圾回收的执行一般会使临时对象池中的对象值被全部移除。
//也就是说，即使我们永远不会显式的从临时对象池取走某一个对象值，该对象值也不会永远待在临时对象池中。它的生命周期取决于垃圾回收任务下一次的执行时间。
//在这里，我们使用runtime/debug代码包的SetGCPercent函数来禁用、恢复GC以及指定垃圾收集比率，以保证我们的演示能够如愿进行。
func main() {
	// 禁用GC，并保证在main函数执行结束前恢复GC。
	defer debug.SetGCPercent(debug.SetGCPercent(-1))
	var count int32
	newFunc := func() interface{} {
		return atomic.AddInt32(&count, 1)
	}
	pool := sync.Pool{New: newFunc}

	// New 字段值的作用。
	v1 := pool.Get()
	fmt.Printf("Value 1: %v\n", v1)

	// 临时对象池的存取。
	pool.Put(10)
	pool.Put(11)
	pool.Put(12)
	v2 := pool.Get()
	fmt.Printf("Value 2: %v\n", v2)

	// 垃圾回收对临时对象池的影响。
	debug.SetGCPercent(100)
	runtime.GC()
	v3 := pool.Get()
	fmt.Printf("Value 3: %v\n", v3)
	pool.New = nil
	v4 := pool.Get()
	fmt.Printf("Value 4: %v\n", v4)
}

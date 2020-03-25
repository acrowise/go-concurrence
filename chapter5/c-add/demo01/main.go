package main

import (
	"fmt"
	"sync/atomic"
)

func main() {

	var i32 int32
	fmt.Printf("i32的的为%v \n",i32)
	//todo 之所以要求第一个参数值必须是指针类型的值，是因为该函数需要获得被操作值在内存中的存放位置，以便施加特殊的CPU指令。
	newi32 := atomic.AddInt32(&i32,3) //给变量i32增大3,不仅返回了新值，同时也更改了原值额
	fmt.Printf("newi32的值为%v,i32的的为%v \n",newi32,i32)

	//原子地 增/减 对应类型的值
	//atomic.AddInt64
	//atomic.AddUint32
	//atomic.AddUint64
	//atomic.AddUintptr
	//注意,并不存在名为 atomic.AddPointer的函数,因为 unsafe.Pointer类型的值无法被加减


	//todo 原子的减少怎么做
	var  i64  int64
	atomic.AddInt64(&i64,-3)
	fmt.Printf("i64的的为%v \n",i64)

    /* 下面这种就不能这样写了,因为uint64类型是大于0的，这样就溢出了，第二个参数的类型就不符合了，怎么办呢？
	  var  iu64  uint64
	  atomic.AddUint64(&iu64,-3)
	  fmt.Printf("iu64的的为%v \n",iu64)
    */

	var  iu64  uint64=10
	atomic.AddUint64(&iu64,^uint64(3-1))
	fmt.Printf("iu64的的为%v \n",iu64)









}

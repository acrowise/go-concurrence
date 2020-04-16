package main

import (
	"fmt"
	"runtime"
)

/*

runtime.GOMAXPROCS(逻辑CPU数量)

这里的逻辑CPU数量可以有如下几种数值：

<1：不修改任何数值。
=1：单核心执行。
>1：多核并发执行。

 */
func main() {
	//查询多少个cpu数量
	num:=runtime.NumCPU()
	fmt.Printf("一共有%d个逻辑cpu\r\n",num)

	// 设置执行使用的核心数
	r:=runtime.GOMAXPROCS(0)
	fmt.Printf("设置执行核心数的返回结果为%d",r)

}

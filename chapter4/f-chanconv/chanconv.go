package main

import "fmt"

func main() {

	//在下面程序中，每一个类型转换表达式的第二个结果值都会是false。
	//因此，利用函数声明将双向通道转换为单向通道的做法，只能算是Go语言的一个语法糖


	var ok bool



	ch := make(chan int, 1)
	_, ok = interface{}(ch).(<-chan int)
	fmt.Println("chan int => <-chan int:", ok)

	_, ok = interface{}(ch).(chan<- int)
	fmt.Println("chan int => chan<- int:", ok)

	//*****************************************************

	sch := make(chan<- int, 1)
	_, ok = interface{}(sch).(chan int)
	fmt.Println("chan<- int => chan int:", ok)

	//**********************************************

	rch := make(<-chan int, 1)
	_, ok = interface{}(rch).(chan int)
	fmt.Println("<-chan int => chan int:", ok)




}

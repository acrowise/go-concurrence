package main
import (
	"fmt"
	"unsafe"
	"reflect"
)
//go语言函数传参可以传递struct，传递pointer，还有传递interface，他们主要区别是什么呢？
//@link https://stackoverflow.com/questions/44370277/type-is-pointer-to-interface-not-interface-confusion


type MyInterface interface {
	test()
}

//@todo 注意 MyStruct是实现了接口MyInterface的类型
type MyStruct struct {
	i1 int64
	i2 int64
	i3 int64
	i4 int64
	i5 int64
}

func (m MyStruct)test()  {

}
//-----------------------------------------------------------------------
//传递结构体类型
func Hello1(p MyStruct) {
	fmt.Println("size:", unsafe.Sizeof(p), "; type:", reflect.TypeOf(p), "; value:", reflect.ValueOf(p))
}
//传递结构体指针类型
func Hello2(p * MyStruct) {
	fmt.Println("size:", unsafe.Sizeof(p), "; type:", reflect.TypeOf(p), "; value:", reflect.ValueOf(p))
}

//接口类型
func Hello3(p MyInterface) {
	fmt.Println("size:", unsafe.Sizeof(p), "; type:", reflect.TypeOf(p), "; value:", reflect.ValueOf(p))
}
//接口指针类型
func Hello4(p * MyInterface) {
  fmt.Println("size:", unsafe.Sizeof(p), "; type:", reflect.TypeOf(p), "; value:", reflect.ValueOf(p))
}

//@todo https://www.jianshu.com/p/42762865c2d8
func main() {
	m := MyStruct { 11,22,33,44,55 }
	Hello1(m)
	Hello2(&m)
	Hello3(m)
	//Hello4((&m)) //这里报错


}
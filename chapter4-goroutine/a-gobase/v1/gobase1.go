package main

//打印G还没有被调度器执行整个程序就嗝屁了
func main() {
	go println("Go! Goroutine!")
}

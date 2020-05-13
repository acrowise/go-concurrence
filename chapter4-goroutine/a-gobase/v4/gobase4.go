package main

import (
	"fmt"
	"time"
)

//todo 输出结果不确定额，大部分是这样的:
/*
Hello, Mark!
Hello, Mark!
Hello, Mark!
Hello, Mark!
Hello, Mark!
 */
//todo 也有可能是这样的额:
/*
Hello, Mark!
Hello, Mark!
Hello, Mark!
Hello, Jim!
Hello, Mark!
 */
func main() {
	names := []string{"Eric", "Harry", "Robert", "Jim", "Mark"}
	for _, name := range names {
		go func() {
			fmt.Printf("Hello, %s!\n", name)
		}()
		//time.Sleep(time.Millisecond)
	}
	time.Sleep(time.Millisecond)
}

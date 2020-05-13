package main

import (
	"fmt"
	"time"
)
//这种写法才能达到我们的效果
func main() {
	names := []string{"Eric", "Harry", "Robert", "Jim", "Mark"}
	for _, name := range names {
		go func(who string) {
			fmt.Printf("Hello, %s!\n", who)
		}(name)
	}
	time.Sleep(time.Millisecond)
}

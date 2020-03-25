
package main

import (
	"fmt"
	"time"
)


//ticker只要定义完成，从此刻开始计时，不需要任何其他的操作，每隔固定时间都会触发。
func main() {

	d := time.Duration(time.Second*2)
	t := time.NewTicker(d)
	defer t.Stop()
	for {
		<- t.C
		fmt.Println("ticker定时任务到期了...")
	}

}




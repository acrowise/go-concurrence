package main

import (
	"fmt"
	"math/rand"
)


//rand.Int31n函数可以在给定的范围内生成一个伪随机数


func main() {

	for i:=0;i<20;i++{
		fmt.Printf("%d,",rand.Int31n(10))  //[0,9]


	}

	fmt.Println()

	for i:=0;i<20;i++{
		fmt.Printf("%d,",rand.Int31n(10)+1)  //[1,10]
	}




	
}

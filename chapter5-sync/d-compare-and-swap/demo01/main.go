package main

import (
	"fmt"
	"sync/atomic"
)

var  value  int32

func main() {
	fmt.Println(value)
	addValue(10)
	fmt.Println(value)
	addValue(2)
	fmt.Println(value)


}


func  addValue(delta int32){
	  for{
	  	   v := value
	  	   if atomic.CompareAndSwapInt32(&value,v,(v + delta)){
	  	   	   break
		   }
	  }
}

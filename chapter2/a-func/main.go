package main

import (
	"errors"
	"fmt"
)

//dividend 被除数
//divisor  除数
func  divide(dividend int,divisor int)(result int,err error){
    if divisor == 0{
    	 err= errors.New("division by zero")
    	 return
	}
    result = dividend / divisor
	return
}

//用于定义二元操作的函数类型
type  binaryOperation func  (operand1 int,operand2 int)(result int,err error)


func  operate(op1 int,op2 int,bop binaryOperation)(result int,err error){
	if bop ==nil{
	   err=errors.New("invalid binary  operation function")
	   return
	}
	return  bop(op1,op2)
}



func main() {

	fmt.Println(divide(10,2))


}

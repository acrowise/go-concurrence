package main

import (
	"fmt"
	"strconv"
)

func main() {
	e := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	for i := 0; i < len(e); i++ {
		ascii, _ := strconv.Atoi(fmt.Sprintf("%d", e[i]))
		fmt.Printf("%c的ascii码值为%d\r\n", e[i], ascii)
	}

}

package main

import (
	"fmt"
	"sync/atomic"
)

//go vet copiedvalue.go

func main() {
	var countVal atomic.Value
	countVal.Store([]int{1, 3, 5, 7})
	anotherStore(countVal)
	fmt.Printf("The count value: %+v \n", countVal.Load()) //The count value: [1 3 5 7]
}

func anotherStore(countVal atomic.Value) {
	countVal.Store([]int{2, 4, 6, 8})
}

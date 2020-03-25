package main

import (
	"bytes"
	"fmt"
)

func main() {

	var buf bytes.Buffer
	buf.WriteString("Checking the parameters...")

	fmt.Println(buf.String())
	
}

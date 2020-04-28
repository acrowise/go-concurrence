package main

import (
	"fmt"
	elasticsearch7 "github.com/elastic/go-elasticsearch/v7"
)

func main(){

	es7, _ := elasticsearch7.NewDefaultClient()
	fmt.Printf("%+v",es7)
}

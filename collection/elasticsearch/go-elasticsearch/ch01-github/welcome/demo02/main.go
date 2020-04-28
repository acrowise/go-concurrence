package main

import (
	"github.com/elastic/go-elasticsearch/v7"
	"log"
)

func main(){


	//使用elasticsearch.NewDefaultClient()函数使用默认设置创建客户端。
	es, _ := elasticsearch.NewDefaultClient()
	log.Println(elasticsearch.Version)
	log.Println(es.Info())
}

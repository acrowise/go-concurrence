package main

import (
	"context"
	"fmt"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"go-collection/elasticsearch/go-elasticsearch/lib/db"
)

//删除索引
func main() {
	req := esapi.IndicesDeleteRequest{
		Index: []string{"ali_index"},
	}
	res, err := req.Do(context.Background(),db.Client)
	db.CheckError(err)
	defer res.Body.Close()
	fmt.Println(res.String())
}

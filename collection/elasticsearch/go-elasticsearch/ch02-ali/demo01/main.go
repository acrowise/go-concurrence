package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"go-collection/elasticsearch/go-elasticsearch/lib/db"
)

//创建索引
func main() {

	body := map[string]interface{}{
		"settings":map[string]interface{}{
			"index":map[string]interface{}{
				"number_of_shards":5,
				"number_of_replicas":1,
			},
		},
		"mappings": map[string]interface{}{
			"properties":map[string]interface{}{
				"post_date":map[string]interface{}{
					    "type":"date",
				},
				"tags":map[string]interface{}{
						"type":"keyword",
				},
				"title":map[string]interface{}{
							"type":"text", // 表示这个字段不分词
					        //"analyzer" : "cjk",
				},
			},
		},
	}
	jsonBody, _ := json.Marshal(body)
	req := esapi.IndicesCreateRequest{
		Index: "ali_index",
		Body:  bytes.NewReader(jsonBody),
	}
	res, err := req.Do(context.Background(),db.Client)
	db.CheckError(err)
	defer res.Body.Close()
	fmt.Println(res.String())


}

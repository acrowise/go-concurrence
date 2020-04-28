package db

import (
	"fmt"
	"github.com/elastic/go-elasticsearch/v7"
	"os"
)

var Client *elasticsearch.Client

func init() {
	var err error
	config := elasticsearch.Config{}
	config.Addresses = []string{"http://127.0.0.1:9200"}
	Client, err = elasticsearch.NewClient(config)
	CheckError(err)
}

func CheckError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
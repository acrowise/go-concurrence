package main
import (
	"crypto/tls"
	"log"
	"net"
	"net/http"
	"time"
	"github.com/elastic/go-elasticsearch/v7"
)

func main() {
	log.SetFlags(0)
	//本示例演示如何配置客户端的传输。注意：这些值仅用于说明目的，不适合用于任何生产用途。 默认传输就足够了。
	cfg := elasticsearch.Config{
		Addresses: []string{"http://localhost:9200"},
		Transport: &http.Transport{
			MaxIdleConnsPerHost:   10,
			ResponseHeaderTimeout: time.Millisecond,
			DialContext:           (&net.Dialer{Timeout: time.Nanosecond}).DialContext,
			TLSClientConfig: &tls.Config{
				MinVersion: tls.VersionTLS11,
				// ...
			},
		},
	}

	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Printf("Error creating the client: %s", err)
	} else {
		log.Println(es.Info())
		// => dial tcp: i/o timeout
	}
}

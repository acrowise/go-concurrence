package main

import (
	"fmt"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
	"log"
)

func main() {
	//viper目前支持的远程配置中心
	//fmt.Printf("%+v",viper.SupportedRemoteProviders)
	v := viper.New()
	v.AddRemoteProvider("etcd", "localhost:2379", "/conf/dev/config.json")

	v.SetConfigType("json") // Need to explicitly set this to json
	if err := v.ReadRemoteConfig(); err != nil {
		log.Println(err)
		return
	}
	fmt.Println("获取配置文件的mysql.url", v.GetString(`mysql.url`))
	fmt.Println("获取配置文件的mysql.username", v.GetString(`mysql.username`))
	fmt.Println("获取配置文件的mysql.password", v.GetString(`mysql.password`))
	fmt.Println("获取配置文件的redis", v.GetStringSlice("redis"))
	fmt.Println("获取配置文件的smtp", v.GetStringMap("smtp"))
}




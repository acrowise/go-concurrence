package main

import (
	"bytes"
	"fmt"

	"github.com/spf13/viper"
)

func main() {
	v := viper.New()
	v.SetConfigType("json") // 设置配置文件的类型
	// 配置文件内容
	var jsonExample = []byte(`
{
  "port": 10666,
  "mysql": {
    "url": "(127.0.0.1:3306)/biezhi",
    "username": "root",
    "password": "123456"
  },
  "redis": ["127.0.0.1:6377", "127.0.0.1:6378", "127.0.0.1:6379"],
  "smtp": {
    "enable": true,
    "addr": "mail_addr",
    "username": "mail_user",
    "password": "mail_password",
    "to": ["xxx@gmail.com", "xxx@163.com"]
  }
}
`)
	//创建io.Reader
	v.ReadConfig(bytes.NewBuffer(jsonExample))

	fmt.Println("获取配置文件的port", v.GetInt("port"))
	fmt.Println("获取配置文件的mysql.url", v.GetString(`mysql.url`))
	fmt.Println("获取配置文件的mysql.username", v.GetString(`mysql.username`))
	fmt.Println("获取配置文件的mysql.password", v.GetString(`mysql.password`))
	fmt.Println("获取配置文件的redis", v.GetStringSlice("redis"))
	fmt.Println("获取配置文件的smtp", v.GetStringMap("smtp"))
}

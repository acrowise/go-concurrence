package main

import (
	"fmt"
	"gopkg.in/ini.v1"
	"log"
)


type App struct {
	PrefixUrl string
	RuntimeRootPath string
	LogSavePath string
	LogSaveName string
	LogFileExt  string
	TimeFormat  string
}


func main() {
	//读取配置文件
	cfg, err := ini.Load("conf/app.ini") //以执行的入口文件的相对位置进行参考
	if err != nil {
		log.Fatalf("setting.Setup, fail to parse 'conf/app.ini': %v", err)
	}

	//映射到对应的结构体
	var AppSetting = &App{}
	err = cfg.Section("app").MapTo(AppSetting)
	if err != nil {
		log.Println(err)
	}

	//输出配置
	fmt.Println(AppSetting)





}

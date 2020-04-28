package main

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

func main() {
	//prefix := "sdk"
	envs := map[string]string{
		"LOG_LEVEL":      "INFO",
		"MODE":           "DEV",
		"MYSQL_USERNAME": "root",
		"MYSQL_PASSWORD": "123456",
	}
	for k, v := range envs {
		//os.Setenv(fmt.Sprintf("%s_%s", prefix, k), v)
		os.Setenv(k, v)
	}

	v := viper.New()
	//v.SetEnvPrefix(prefix)
	v.AutomaticEnv()

	for k, _ := range envs {
		fmt.Printf("env `%s` = %s\n", k, v.GetString(k))
	}
}

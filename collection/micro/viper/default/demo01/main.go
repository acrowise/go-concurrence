package main

import (
	"fmt"
	"github.com/spf13/viper"
)

func main() {

	// set default config
	viper.SetDefault("ContentDir", "content")
	viper.SetDefault("LayoutDir", "layouts")
	viper.SetDefault("Taxonomies", map[string]string{"name": "jack", "sex": "male"})

	fmt.Println(viper.GetBool("ContentDir"))
	fmt.Println(viper.GetString("LayoutDir"))
	fmt.Println(viper.GetStringMapString("Taxonomies"))

	
}

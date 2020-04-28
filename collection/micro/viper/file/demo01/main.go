package main
import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"log"
	"time"
)


//viper.AddConfigPath("/etc/appname/")   //设置配置文件的搜索目录
//viper.AddConfigPath("$HOME/.appname")  // 设置配置文件的搜索目录

func init() {
	//1.设置配置文件名为 config, 不需要配置文件扩展名，配置文件的类型 viper 会自动根据扩展名自动匹配.
	viper.SetConfigName("config") // 读取json配置文件
	//viper.SetConfigName("y") // 读取yaml配置文件
	//2.设置配置文件搜索的目录，. 表示和当前编译好的二进制文件在同一个目录。可以添加多个配置文件目录，如在第一个目录中找到就不不继续到其他目录中查找.
	viper.AddConfigPath(".")
	//3.加载配置文件内容
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Println("no such config file")
		} else {
			log.Println("read config error")
		}
		log.Fatal(err) // 读取配置文件失败致命错误
	}
}

func main() {
	//4.获取配置文件中配置项的信息
	fmt.Println("获取配置文件的port", viper.GetInt("port"))
	fmt.Println("获取配置文件的mysql.url", viper.GetString(`mysql.url`))
	fmt.Println("获取配置文件的mysql.username", viper.GetString(`mysql.username`))
	fmt.Println("获取配置文件的mysql.password", viper.GetString(`mysql.password`))
	fmt.Println("获取配置文件的redis", viper.GetStringSlice("redis"))
	fmt.Println("获取配置文件的smtp", viper.GetStringMap("smtp"))


	//5.监控配置变化情况
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		//viper配置发生变化了 执行响应的操作
		fmt.Println("Config file changed:", e.Name)
	})


	time.Sleep(1000 * time.Second)


}


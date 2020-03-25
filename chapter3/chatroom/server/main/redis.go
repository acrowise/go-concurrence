package main
import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"time"
)

//定义一个全局的pool
var pool *redis.Pool

func initPool(address string, maxIdle, maxActive int, idleTimeout time.Duration) {

	pool = &redis.Pool{
		MaxIdle: maxIdle, //最大空闲链接数
		MaxActive: maxActive, // 表示和数据库的最大链接数， 0 表示没有限制
		IdleTimeout: idleTimeout, // 最大空闲时间
		Dial: func() (redis.Conn, error) { // 初始化链接的代码， 链接哪个ip的redis
		   //return redis.Dial("tcp", address)
			c, err := redis.Dial("tcp", address)
			if err != nil {
				fmt.Println("Dial err=",err)
				return nil, err
			}
			if _, err := c.Do("AUTH","root"); err != nil {
				c.Close()
				fmt.Println("AUTH err=",err)
				return nil, err
			}
			if _, err := c.Do("SELECT",0); err != nil {
				c.Close()
				fmt.Println("SELECT err=",err)
				return nil, err
			}
			return c, nil

		},
	}
}
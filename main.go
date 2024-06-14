package main

import "github.com/gin-gonic/gin"

func main() {
	// 用默认配置创建实例
	r := gin.Default() //里面有log日志或者recovery持久化进程等中间件
	// 注册路由
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong jserWang",
		})
	})
	// 绑定端口和ip，运行实例
	r.Run("localhost:8000") // listen and serve on 0.0.0.0:8080
}

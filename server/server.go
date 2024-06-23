package server

import (
	"RealWorldWeb/handlers"
	"github.com/gin-gonic/gin"
)

func RunHttpServer() {
	// 用默认配置创建实例
	r := gin.Default() //里面有log日志或者recovery持久化进程等中间件
	// 注册路由
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong jserWang",
		})
	})
	// 绑定端口和ip，运行实例

	handlers.AddUserHandler(r)
	handlers.AddArticleHandler(r)
	handlers.AddTagsHandler(r)
	handlers.AddArticleCommentHandler(r)

	r.Run("localhost:8000") // listen and serve on 0.0.0.0:8080
}

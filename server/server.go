package server

import (
	"RealWorldWeb/handlers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"time"
)

func RunHttpServer() {
	// 用默认配置创建实例
	r := gin.Default() //里面有log日志或者recovery持久化进程等中间件
	// 绑定端口和ip，运行实例
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"*"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	handlers.AddUserHandler(r)
	handlers.AddArticleHandler(r)
	handlers.AddTagsHandler(r)
	handlers.AddArticleCommentHandler(r)

	r.Run(":8000") // listen and serve on 0.0.0.0:8080
}

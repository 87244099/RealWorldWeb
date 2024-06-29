package handlers

import (
	"RealWorldWeb/logger"
	"github.com/gin-gonic/gin"
	"io"
	"time"
)

func AddSSE(r *gin.Engine) {
	sseGroup := r.Group("/api/sse")
	sseGroup.GET("", getSse)
}

func getSse(ctx *gin.Context) {
	log := logger.New(ctx)
	log.Infof("get sse")
	ticker := time.NewTicker(1 * time.Second)
	ctx.Stream(func(w io.Writer) bool {
		select {
		case <-ticker.C:
			log.Infof("send sse event")
			ctx.SSEvent("", "heartbeat: "+time.Now().String())
		}
		return true
	})

	/*
		前端控制台可以输入如下代码:
		let es = new EventSource("/api/sse")
		es.onmessage = function(event){
			console.log(event.data);
		}

		es.close() 主动结束
	*/
}

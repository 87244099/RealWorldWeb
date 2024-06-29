package handlers

import (
	"RealWorldWeb/logger"
	"RealWorldWeb/utils"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"time"
)

var upgrader = websocket.Upgrader{}

func AddWebsocket(r *gin.Engine) {
	wsGroup := r.Group("/api/ws")
	wsGroup.GET("", ws)
}

func ws(ctx *gin.Context) {
	log := logger.New(ctx)
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.WithError(err).Errorf("Upgrader ws failed")
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	defer conn.Close()
	//开另外一条协程定时发送心跳
	go func() {
		for {
			time.Sleep(2 * time.Second)
			conn.WriteJSON(map[string]interface{}{
				"type": "heartbeat",
			})
		}
	}()

	for {
		reqMsg := make(map[string]interface{})
		err := conn.ReadJSON(&reqMsg)
		if err != nil {
			log.WithError(err).Errorf("Read JSON failed")
			return
		}
		log.Infof("read msg: %v\n", utils.JsonMarshal(reqMsg))

		if reqMsg["exit"] != nil {
			return
		}

		err = conn.WriteJSON(reqMsg)
		if err != nil {
			log.WithError(err).Errorf("WriteJSON failed")
			return
		}
	}
	/**
	客户端测试代码：https://en.wikipedia.org/wiki/WebSocket

	// Connect to server
	ws = new WebSocket("ws://127.0.0.1/scoreboard") // Local server
	// ws = new WebSocket("wss://game.example.com/scoreboard") // Remote server

	ws.onopen = () => {
	    console.log("Connection opened")
	    ws.send("Hi server, please send me the score of yesterday's game")
	}

	ws.onmessage = (event) => {
	    console.log("Data received", event.data)
	    ws.close() // We got the score so we don't need the connection anymore
	}

	ws.onclose = (event) => {
	    console.log("Connection closed", event.code, event.reason, event.wasClean)
	}

	ws.onerror = () => {
	    console.log("Connection closed due to error")
	}
	*/
}

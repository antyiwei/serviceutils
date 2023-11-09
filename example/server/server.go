package main

import (
	"log"

	"github.com/antyiwei/serviceutils/server"
	"github.com/antyiwei/serviceutils/wsserver"
	"github.com/gin-gonic/gin"
)

type TestReq struct {
	Data struct {
		Word string
	}
}

func main() {

	server.StartServer(":9999", func(engine *gin.Engine, ws *wsserver.WsEngine) {

		ws.On("some", func(wc *wsserver.WsRequest, c *gin.Context) {

			log.Println("hi")
			var req TestReq

			wc.Bind(&req)

			wc.ReturnSuccess("hi")
		})
	})
}

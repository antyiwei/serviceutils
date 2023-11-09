package server

import (
	"github.com/antyiwei/serviceutils/wsserver"
	"github.com/gin-gonic/gin"
	"github.com/luaxlou/goutils/gindefault"
)

func StartServer(addr string, initial func(engine *gin.Engine, ws *wsserver.WsEngine)) {
	gindefault.Run(addr, func(engine *gin.Engine) {

		engine.GET("/healthz", func(c *gin.Context) {

			c.String(200, "ok")

		})

		wsEngine := wsserver.NewServer()

		engine.GET("/ws", wsserver.WsReadHandler(wsEngine))

		initial(engine, wsEngine)

	})
}

package server

import (
	"github.com/antyiwei/serviceutils/server/gindefault"
	"github.com/antyiwei/serviceutils/wsserver"
	"github.com/gin-gonic/gin"
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

func StartServer2(addr string, initial func(engine *gin.Engine, ws *wsserver.WsEngine), doSomethingBeforeExit func() error) {
	gindefault.Run2(addr, func(engine *gin.Engine) {

		engine.GET("/healthz", func(c *gin.Context) {

			c.String(200, "ok")

		})

		wsEngine := wsserver.NewServer()

		engine.GET("/ws", wsserver.WsReadHandler(wsEngine))

		initial(engine, wsEngine)

	}, doSomethingBeforeExit)
}

package wsserver

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/antyiwei/goutils/debugutils"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type WsEngine struct {
	authFunc   func(c *gin.Context) error
	WsHandlers map[string]func(ws *WsRequest, c *gin.Context)
}

func NewServer() *WsEngine {

	return &WsEngine{WsHandlers: make(map[string]func(ws *WsRequest, c *gin.Context), 0)}
}

func (we *WsEngine) Auth(authFunc func(c *gin.Context) error) {

	we.authFunc = authFunc
}

func (we *WsEngine) On(command string, handler func(wc *WsRequest, c *gin.Context)) {

	we.WsHandlers[command] = handler
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func WsReadHandler(we *WsEngine) func(c *gin.Context) {

	return func(c *gin.Context) {

		if we.authFunc != nil {

			if err := we.authFunc(c); err != nil {
				c.String(403, err.Error())
				c.Abort()

				return
			}

		}

		log.Printf("new wsclient connected: %s", c.ClientIP())
		ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)

		if err != nil {
			log.Print("upgrade:", err)
			return
		}

		var mux sync.RWMutex

		ws.SetPingHandler(func(appData string) error {
			mux.Lock()

			ws.WriteMessage(websocket.PongMessage, []byte{})
			mux.Unlock()
			return nil
		})

		defer ws.Close()
		for {
			mt, message, err := ws.ReadMessage()

			if err != nil {
				debugutils.Println("read:", err.Error())
				break
			}

			switch mt {

			case websocket.BinaryMessage:

				var tc TempCommand

				err := json.Unmarshal(message, &tc)

				if err != nil {

					ret := err.Error()
					mux.Lock()

					ws.WriteMessage(websocket.TextMessage, []byte(ret))
					mux.Unlock()

				} else {
					handler, ok := we.WsHandlers[tc.Command]

					context := &WsRequest{RequestId: tc.RequestId,
						Command: tc.Command,
						Mt:      mt,
						Data:    message,
						Ws:      ws,
						Mux:     mux,
					}

					if !ok {
						context.ReturnFail(404, "not support command")

					} else {
						handler(context, c)
					}

				}
			}

		}
	}
}

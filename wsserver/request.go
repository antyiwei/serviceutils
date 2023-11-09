package wsserver

import (
	"encoding/json"
	"errors"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type WsRequest struct {
	MerchantId int64
	RequestId  string
	Command    string
	Mt         int
	Data       []byte

	Mux sync.RWMutex

	Ws *websocket.Conn
	C  *gin.Context
}

type Message struct {
	MessageType int
	Data        []byte
}

func (wc *WsRequest) GetResult() (r *WsResult) {

	if wc.Mt != websocket.BinaryMessage {
		return nil
	}

	json.Unmarshal(wc.Data, &r)

	return

}

func (wc *WsRequest) Bind(bind interface{}) error {

	if wc.Mt != websocket.BinaryMessage {
		return errors.New("not binary type message")
	}

	return json.Unmarshal(wc.Data, bind)

}

func (wc *WsRequest) GetString() string {

	if wc.Mt != websocket.BinaryMessage {
		return ""
	}

	var cmd WsCommand

	err := json.Unmarshal(wc.Data, &cmd)

	if err != nil {
		return ""
	}

	str, ok := cmd.Data.(string)

	if !ok {
		return ""
	}

	return str

}

func (wc *WsRequest) String(msg string) {


	wc.WriteMessage(websocket.TextMessage, []byte(msg))

}

func (wc *WsRequest) ReturnSuccess(msg string) {

	res := WsResult{
		RequestId: wc.RequestId,
		Type:      1,

		Command: wc.Command,

		Status: 200,
		Msg:    msg,
	}

	bs, _ := json.Marshal(&res)

	wc.WriteMessage(websocket.BinaryMessage, bs)

}


func (wc *WsRequest) WriteMessage(messageType int ,msg []byte) {

	wc.Mux.Lock()
	wc.Ws.WriteMessage(messageType, msg)
	wc.Mux.Unlock()
}


func (wc *WsRequest) ReturnData(data interface{}) {

	res := WsResult{
		RequestId: wc.RequestId,
		Type:      1,

		Command: wc.Command,

		Status: 200,
		Data:   data,
	}

	bs, _ := json.Marshal(&res)

	wc.WriteMessage(websocket.BinaryMessage, bs)

}

func (wc *WsRequest) ReturnError(err error) {

	wc.ReturnFail(500, err.Error())
}

func (wc *WsRequest) ReturnFail(code int, msg string) {

	res := WsResult{
		RequestId: wc.RequestId,

		Type: 1,

		Status: code,
		Msg:    msg,
	}

	bs, _ := json.Marshal(&res)
	wc.WriteMessage(websocket.BinaryMessage, bs)

}

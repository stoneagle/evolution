package resp

import (
	"encoding/json"
	"evolution/backend/common/logger"
	"evolution/backend/common/utils"

	melody "gopkg.in/olahol/melody.v1"
)

var (
	wsIDKey string = "wsID"
	wsCBKey string = "wsCB"
)

type Websocket struct {
	Intence *melody.Melody
}

type WsCallback func([]byte) WebsocketResponse

type WebsocketResponse struct {
	Code   ErrorCode   `json:"code"`
	Status WsStatus    `json:"status"`
	Data   interface{} `json:"data"`
	Desc   string      `json:"desc"`
}

type WebsocketContext struct {
}

func NewWebsocket() *Websocket {
	ws := &Websocket{
		Intence: melody.New(),
	}
	ws.Intence.HandleMessage(ws.HandleMessage)
	ws.Intence.HandleConnect(ws.HandleConnect)
	ws.Intence.HandleDisconnect(ws.HandleDisconnect)
	return ws
}

func (ws *Websocket) BuildContext(callback WsCallback) map[string]interface{} {
	context := make(map[string]interface{})
	context[wsIDKey] = utils.UniqueId()
	context[wsCBKey] = callback
	return context
}

func (ws *Websocket) HandleMessage(s *melody.Session, msg []byte) {
	var res WebsocketResponse
	ctx := s.Keys

	logger.Get().Infow("websocket-request:【" + string(msg) + "】")

	res = ctx[wsCBKey].(WsCallback)(msg)
	ret, _ := json.Marshal(&res)

	logger.Get().Infow("websocket-response:【" + string(ret) + "】")

	s.Write(ret)
}

func (ws *Websocket) HandleConnect(s *melody.Session) {
	res := WebsocketResponse{
		Status: WsConnect,
		Code:   ErrorOk,
		Data:   s.Keys[wsIDKey],
		Desc:   "connect",
	}
	ret, _ := json.Marshal(&res)
	s.Write(ret)
}

func (ws *Websocket) HandleDisconnect(s *melody.Session) {
}

func (ws *Websocket) Message(data interface{}) WebsocketResponse {
	return WebsocketResponse{
		Status: WsMessage,
		Code:   ErrorOk,
		Data:   data,
		Desc:   "success",
	}
}

func (ws *Websocket) BusinessError(errorCode ErrorCode, desc string, err error) WebsocketResponse {
	if err != nil {
		desc += ":" + err.Error()
	}
	return WebsocketResponse{
		Status: WsMessage,
		Code:   errorCode,
		Data:   struct{}{},
		Desc:   desc,
	}
}

func (ws *Websocket) ServerError(errorCode ErrorCode, desc string, err error) WebsocketResponse {
	if err != nil {
		desc += ":" + err.Error()
	}
	return WebsocketResponse{
		Status: WsError,
		Code:   errorCode,
		Data:   struct{}{},
		Desc:   desc,
	}
}

package transport

import (
	"github.com/autogrowsystems/faye-go/protocol"
	"github.com/autogrowsystems/faye-go/utils"
	"github.com/gorilla/websocket"
	"io"
)

const WebSocketConnectionPriority = 10

type Server interface {
	HandleRequest(interface{}, protocol.Connection)
	Logger() utils.Logger
}

type WebSocketConnection struct {
	ws         *websocket.Conn
	failedSend bool
}

func (wc *WebSocketConnection) Send(msgs []protocol.Message) error {
	err := wc.ws.WriteJSON(msgs)
	if err != nil {
		wc.failedSend = true
	}
	return err
}

func (wc *WebSocketConnection) IsConnected() bool {
	return wc.failedSend
}

func (wc *WebSocketConnection) Close() {
	wc.ws.Close()
}

func (wc WebSocketConnection) Priority() int {
	return WebSocketConnectionPriority
}

func (lp WebSocketConnection) IsSingleShot() bool {
	return false
}

func WebsocketServer(m Server) func(*websocket.Conn) {
	return func(ws *websocket.Conn) {
		var data interface{}
		wsConn := WebSocketConnection{ws, true}
		for {
			err := ws.ReadJSON(&data)
			if err != nil {
				if err == io.EOF {
					m.Logger().Errorf("EOF while reading from socket")
					return
				} else {
					m.Logger().Errorf("While reading from socket: %s", err)
					return
				}
			}

			if arr := data.([]interface{}); len(arr) == 0 {
				ws.WriteJSON([]string{})
			} else {
				m.HandleRequest(data, &wsConn)
			}
		}
	}
}

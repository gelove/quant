package socket

import (
	"log"
	"time"
)

// WsHandler handle raw websocket message
type WsHandler func(*Socket, []byte)

// var locker = &sync.Mutex{}

// func Reconnect(ws *Socket, err error) {
// 	locker.Lock()
// 	defer locker.Unlock()
// 	log.Println("Reconnect IsConnected => ", ws.IsConnected())
// 	if !ws.IsConnected() {
// 		ws.Connect()
// 	}
// }

func Run(url string, handler WsHandler) (ws *Socket) {
	ws = New(url, SetIsKeepAlive(true), SetInterval(30*time.Second))

	ws.OnConnected = func(ws *Socket) {
		log.Println("OnConnected success")
	}
	ws.OnTextMessage = handler
	ws.OnDisconnected = Reconnect
	ws.OnConnectError = Reconnect
	ws.Connect()
	return
}

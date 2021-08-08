package test

import (
	"os"
	"os/signal"
	"quant/pkg/socket"
	"testing"
	"time"
)

func TestReconnect(t *testing.T) {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	connectionOptions := socket.ConnectionOptions{
		//Proxy: BuildProxy("http://127.0.0.1:1087"),
		UseSSL:         false,
		UseCompression: false,
		Subprotocols:   []string{"chat", "superchat"},
	}
	url := "ws://echo.websocket.org/"
	ws := socket.New(url, socket.SetConnectionOptions(connectionOptions), socket.SetInterval(30*time.Second))
	ws.SetHeader("Accept-Encoding", "gzip, deflate, sdch")
	ws.SetHeader("Accept-Language", "en-US,en;q=0.8")
	ws.SetHeader("Pragma", "no-cache")
	ws.SetHeader("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/49.0.2623.87 Safari/537.36")

	ws.OnConnected = func(socket *socket.Socket) {
		t.Log("OnConnected success")
	}
	ws.OnTextMessage = func(socket *socket.Socket, message []byte) {
		t.Logf("Received message %s", message)
	}
	ws.OnPingReceived = func(socket *socket.Socket, data string) {
		t.Logf("Received ping %s", data)
	}
	ws.OnDisconnected = socket.Reconnect
	ws.OnConnectError = socket.Reconnect

	ws.Connect()

	go func() {
		for {
			time.Sleep(1000 * time.Millisecond)
			ws.SendText("This is my sample test message")
		}
	}()

	testReconnect := time.After(8 * time.Second)
	testFinished := time.After(24 * time.Second)
	for {
		select {
		case <-testReconnect:
			t.Log("Connect timeout")
			ws.Close(false)
			continue
		case <-testFinished:
			t.Log("Connect finished")
			ws.Close(true)
			return
		case <-interrupt:
			t.Log("interrupt")
			ws.Close(true)
			return
		}
	}
}

package socket

import (
	"crypto/tls"
	"errors"
	"log"
	"net/http"
	"net/url"
	"quant/pkg/utils/logger"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
)

const (
	StateDisconnected int64 = iota
	StateConnected
)

type Socket struct {
	url               string
	isKeepAlive       bool
	state             int64
	conn              *websocket.Conn
	websocketDialer   *websocket.Dialer
	connectionOptions ConnectionOptions
	requestHeader     http.Header
	OnConnected       func(socket *Socket)
	OnTextMessage     func(socket *Socket, message []byte)
	OnBinaryMessage   func(socket *Socket, data []byte)
	OnConnectError    func(socket *Socket, err error)
	OnDisconnected    func(socket *Socket, err error)
	OnPingReceived    func(socket *Socket, data string)
	OnPongReceived    func(socket *Socket, data string)
	timeout           time.Duration
	interval          time.Duration
	connectMu         sync.RWMutex // 连接句柄锁
	reconnectMu       sync.RWMutex // 断线重连锁
	sendMu            sync.Mutex   // 在另一个goroutine中定时写入pong来保持连接, 多个goroutine会导致并发写入
	// receiveMu         sync.Mutex   // 本程序只有一个goroutine读, 应该可以不用读锁
}

type ConnectionOptions struct {
	UseCompression bool
	UseSSL         bool
	Proxy          func(*http.Request) (*url.URL, error)
	Subprotocols   []string
}

// SetOption 定义配置项
type SetOption func(*Socket)

func SetIsKeepAlive(isKeepAlive bool) SetOption {
	return func(s *Socket) {
		s.isKeepAlive = isKeepAlive
	}
}

func SetInterval(interval time.Duration) SetOption {
	return func(s *Socket) {
		s.interval = interval
	}
}

func SetTimeout(timeout time.Duration) SetOption {
	return func(s *Socket) {
		s.timeout = timeout
	}
}

func SetConnectionOptions(options ConnectionOptions) SetOption {
	return func(s *Socket) {
		s.connectionOptions = options
	}
}

func New(url string, opts ...SetOption) *Socket {
	ws := &Socket{
		url:             url,
		requestHeader:   http.Header{},
		websocketDialer: &websocket.Dialer{},
		interval:        60 * time.Second,
		connectionOptions: ConnectionOptions{
			UseCompression: false,
			UseSSL:         true,
		},
	}
	for _, opt := range opts {
		opt(ws)
	}
	return ws
}

func (ws *Socket) Url() string {
	return ws.url
}

func (ws *Socket) getConn() *websocket.Conn {
	ws.connectMu.RLock()
	defer ws.connectMu.RUnlock()
	return ws.conn
}

func (ws *Socket) setConn(val *websocket.Conn) {
	ws.connectMu.Lock()
	defer ws.connectMu.Unlock()
	ws.conn = val
}

func (ws *Socket) IsKeepAlive() bool {
	return ws.isKeepAlive
}

func (ws *Socket) Interval() time.Duration {
	return ws.interval
}

func (ws *Socket) Timeout() time.Duration {
	return ws.timeout
}

func (ws *Socket) setConnectionConfig() {
	ws.websocketDialer.EnableCompression = ws.connectionOptions.UseCompression
	ws.websocketDialer.TLSClientConfig = &tls.Config{InsecureSkipVerify: ws.connectionOptions.UseSSL}
	ws.websocketDialer.Proxy = ws.connectionOptions.Proxy
	ws.websocketDialer.Subprotocols = ws.connectionOptions.Subprotocols
}

func (ws *Socket) SetHeader(key, value string) {
	ws.requestHeader.Set(key, value)
}

func (ws *Socket) Connect() {
	ws.setConnectionConfig()

	conn, resp, err := ws.websocketDialer.Dial(ws.url, ws.requestHeader)

	if err != nil {
		logger.Errorf("client Connect to server err: %+v", err)
		if resp != nil {
			logger.Errorf("HTTP error response %d status: %s", resp.StatusCode, resp.Status)
		}
		atomic.StoreInt64(&ws.state, StateDisconnected)
		if ws.OnConnectError != nil {
			ws.OnConnectError(ws, err)
		}
		return
	}

	ws.setConn(conn)
	atomic.StoreInt64(&ws.state, StateConnected)
	log.Println("Connected to server")

	if ws.OnConnected != nil {
		ws.OnConnected(ws)
	}

	defaultPingHandler := ws.getConn().PingHandler()
	ws.getConn().SetPingHandler(func(data string) error {
		logger.Info("Received PING from server")
		if ws.OnPingReceived != nil {
			ws.OnPingReceived(ws, data)
		}
		return defaultPingHandler(data)
	})

	defaultPongHandler := ws.getConn().PongHandler()
	ws.getConn().SetPongHandler(func(data string) error {
		logger.Info("Received PONG from server")
		if ws.OnPongReceived != nil {
			ws.OnPongReceived(ws, data)
		}
		return defaultPongHandler(data)
	})

	defaultCloseHandler := ws.getConn().CloseHandler()
	ws.getConn().SetCloseHandler(func(code int, text string) error {
		err := defaultCloseHandler(code, text)
		if err != nil {
			logger.Errorf("Disconnected from server err: %+v", err)
		}
		atomic.StoreInt64(&ws.state, StateDisconnected)
		if ws.OnDisconnected != nil {
			logger.Infof("CloseHandler exec OnDisconnected, err: %+v", err)
			ws.OnDisconnected(ws, errors.New(text))
		}
		return err
	})

	go func() {
		for {
			if ws.timeout != 0 {
				ws.getConn().SetReadDeadline(time.Now().Add(ws.timeout))
			}
			messageType, message, err := ws.receive()
			if err != nil {
				return
			}

			switch messageType {
			case websocket.TextMessage:
				if ws.OnTextMessage != nil {
					ws.OnTextMessage(ws, message)
				}
			case websocket.BinaryMessage:
				if ws.OnBinaryMessage != nil {
					ws.OnBinaryMessage(ws, message)
				}
			}
		}
	}()

	if !ws.isKeepAlive {
		return
	}

	go func() {
		ws.keepAlive()
	}()
}

func (ws *Socket) SendText(message string) {
	err := ws.send(websocket.TextMessage, []byte(message))
	if err != nil {
		logger.Errorf("SendText err: %+v", err)
		return
	}
}

func (ws *Socket) SendBinary(data []byte) {
	err := ws.send(websocket.BinaryMessage, data)
	if err != nil {
		logger.Errorf("SendBinary err: %+v", err)
		return
	}
}

func (ws *Socket) send(messageType int, data []byte) error {
	ws.sendMu.Lock()
	defer ws.sendMu.Unlock()
	err := ws.getConn().WriteMessage(messageType, data)
	return err
}

func (ws *Socket) receive() (messageType int, message []byte, err error) {
	// ws.receiveMu.Lock()
	// defer ws.receiveMu.Unlock()
	messageType, message, err = ws.getConn().ReadMessage()
	if err != nil {
		logger.Errorf("client receive err: %+v", err)
		atomic.StoreInt64(&ws.state, StateDisconnected)
		if ws.OnConnectError != nil {
			ws.OnConnectError(ws, err)
		}
		return
	}
	log.Printf("receive: %s", message)
	return
}

func (ws *Socket) Close(force bool) {
	err := ws.send(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		logger.Errorf("client write CloseMessage err: %+v", err)
	}
	err = ws.getConn().Close()
	if err != nil {
		logger.Errorf("client close err: %+v", err)
	}
	atomic.StoreInt64(&ws.state, StateDisconnected)
	if ws.OnDisconnected != nil && !force {
		logger.Infof("Socket Close exec OnDisconnected, err: %+v", err)
		ws.OnDisconnected(ws, err)
	}
}

func BuildProxy(Url string) func(*http.Request) (*url.URL, error) {
	uProxy, err := url.Parse(Url)
	if err != nil {
		log.Fatal("Error while parsing url ", err)
	}
	return http.ProxyURL(uProxy)
}

func (ws *Socket) keepAlive() {
	ticker := time.NewTicker(ws.interval)
	for range ticker.C {
		ws.sendPong()
	}
}

func (ws *Socket) sendPong() {
	ws.sendMu.Lock()
	defer ws.sendMu.Unlock()
	deadline := time.Now().Add(10 * time.Second)
	err := ws.getConn().WriteControl(websocket.PongMessage, []byte{}, deadline)
	if err != nil {
		logger.Errorf("sendPong err: %+v", err)
		atomic.StoreInt64(&ws.state, StateDisconnected)
		if ws.OnConnectError != nil {
			ws.OnConnectError(ws, err)
		}
		return
	}
	logger.Info("sendPong success")
}

func Reconnect(ws *Socket, err error) {
	if atomic.LoadInt64(&ws.state) == StateConnected {
		return
	}
	ws.reconnectMu.Lock()
	defer ws.reconnectMu.Unlock()
	if atomic.LoadInt64(&ws.state) == StateConnected {
		return
	}
	logger.Infof("Reconnect success, state => %d", atomic.LoadInt64(&ws.state))
	ws.Connect()
}

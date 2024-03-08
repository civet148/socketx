package websock

import (
	"crypto/tls"
	"fmt"
	"github.com/civet148/gotools/parser"
	"github.com/civet148/log"
	"github.com/civet148/socketx/api"
	"github.com/civet148/socketx/types"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
)

type socket struct {
	ui        *parser.UrlInfo
	conn      *websocket.Conn
	accepting chan *websocket.Conn
	closed    bool
	locker    sync.RWMutex
	option    *api.SocketOption
}

func init() {
	_ = api.Register(types.SocketType_WEB, NewSocket)
}

func NewSocket(ui *parser.UrlInfo, options ...api.SocketOption) api.Socket {
	var option *api.SocketOption
	if len(options) != 0 {
		option = &options[0]
	}
	return &socket{
		ui:        ui,
		option:    option,
		accepting: make(chan *websocket.Conn, 1000),
	}
}

func (s *socket) Listen() (err error) {
	engine := gin.Default()
	if s.ui.GetPath() == "" {
		s.ui.Path = "/"
	}

	engine.GET(s.ui.Path, s.webSocketRegister)
	strCertFile := s.ui.Queries[types.WSS_TLS_CERT]
	strKeyFile := s.ui.Queries[types.WSS_TLS_KEY]

	go func() {
		if s.ui.Scheme == types.URL_SCHEME_WSS {
			err = engine.RunTLS(s.ui.Host, strCertFile, strKeyFile)
		} else {
			err = engine.Run(s.ui.Host)
		}

		if err != nil {
			s.closed = true
			log.Errorf("listen websocket closing with error [%v]", err.Error())
			return
		}
	}()
	return
}

func (s *socket) Accept() api.Socket {

	var c *websocket.Conn
	select {
	case c = <-s.accepting:
		{
			c.SetCloseHandler(s.webSocketCloseHandler)
			c.SetPingHandler(s.websocketPingHandler)
			c.SetPongHandler(s.websocketPongHandler)
			return &socket{
				conn: c,
				ui:   s.ui,
			}
		}
	}
	return nil
}

func (s *socket) Connect() (err error) {
	url := fmt.Sprintf("%v://%v%v", s.ui.Scheme, s.ui.Host, s.ui.Path)
	dialer := &websocket.Dialer{}
	if s.ui.Scheme == types.URL_SCHEME_WSS {
		dialer.TLSClientConfig = &tls.Config{RootCAs: nil, InsecureSkipVerify: true}
	}
	var header http.Header
	if s.option != nil {
		header = s.option.Header
	}
	log.Infof("connecting to [%s] with header [%+v]", url, header)
	if s.conn, _, err = dialer.Dial(url, header); err != nil {
		log.Errorf(err.Error())
		return
	}
	return
}

func (s *socket) Send(data []byte, to ...string) (n int, err error) {
	if s.conn == nil {
		err = fmt.Errorf("web socket connection is nil")
		return
	}
	s.locker.Lock()
	defer s.locker.Unlock()
	if err = s.conn.WriteMessage(websocket.BinaryMessage, data); err != nil {
		return
	}
	n = len(data)
	return
}

func (s *socket) Recv(length int) (msg *api.SockMessage, err error) {
	if s.conn == nil {
		return nil, log.Errorf("web socket connection is nil")
	}
	var msgType int
	var data []byte
	if msgType, data, err = s.conn.ReadMessage(); err != nil {
		log.Errorf(err.Error())
		return
	}
	from := s.conn.RemoteAddr().String()
	return &api.SockMessage{
		Sock:    s,
		Data:    data,
		From:    from,
		MsgType: msgType,
	}, nil
}

func (s *socket) Close() (err error) {

	if s.closed {
		return fmt.Errorf("socket already closed")
	}
	s.closed = true
	if s.conn == nil {
		return fmt.Errorf("socket is nil")
	}
	s.closed = true
	return s.conn.Close()
}

func (s *socket) GetLocalAddr() (addr string) {
	if s.conn == nil {
		return s.ui.Host //web socket server connection is nil
	}
	addr = s.conn.LocalAddr().String()
	return
}

func (s *socket) GetRemoteAddr() (addr string) {
	if s.conn == nil {
		return //web socket client connection can't be nil
	}
	addr = s.conn.RemoteAddr().String()
	return
}

func (s *socket) GetSocketType() types.SocketType {
	return types.SocketType_WEB
}

func (s *socket) debugMessageType(msgType int) {

	switch msgType {
	case websocket.TextMessage:
		log.Debugf("message type [TextMessage]")
	case websocket.BinaryMessage:
		log.Debugf("message type [BinaryMessage]")
	case websocket.CloseMessage:
		log.Debugf("message type [CloseMessage]")
	case websocket.PingMessage:
		log.Debugf("message type [PingMessage]")
	case websocket.PongMessage:
		log.Debugf("message type [PongMessage]")
	}
}

func (s *socket) webSocketRegister(ctx *gin.Context) {
	var err error
	upGrader := &websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		Subprotocols: []string{ctx.GetHeader("Sec-WebSocket-Protocol")},
	}
	var c *websocket.Conn
	if c, err = upGrader.Upgrade(ctx.Writer, ctx.Request, nil); err != nil {
		log.Errorf(err.Error())
		return
	}
	//log.Debugf("client [%v] registered", c.RemoteAddr().String())
	s.accepting <- c
}

func (s *socket) webSocketCloseHandler(code int, text string) (err error) {
	log.Debugf("close code [%v] text [%v]", code, text)
	return
}

func (s *socket) websocketPingHandler(appData string) (err error) {
	log.Debugf("ping app data [%v]", appData)
	return
}

func (s *socket) websocketPongHandler(appData string) (err error) {
	log.Debugf("pong app data [%v]", appData)
	return
}

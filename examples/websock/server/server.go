package main

import (
	"github.com/civet148/log"
	"github.com/civet148/socketx"
)

const (
	//WEBSOCKET_SERVER_URL = "ws://0.0.0.0:6668/websocket"
	WEBSOCKET_SERVER_URL = "wss://0.0.0.0:6668/websocket?cert=cert.pem&key=key.pem"
)

const (
	WEBSOCKET_DATA_PING = "ping"
	WEBSOCKET_DATA_PONG = "pong"
)

type ServerHandler struct {
}

func main() {
	log.SetLevel("debug")
	var handler ServerHandler
	sock := socketx.NewServer(WEBSOCKET_SERVER_URL)
	_ = sock.Listen(&handler)
}

func (s *ServerHandler) OnAccept(c *socketx.SocketClient) {
	log.Infof("connection accepted [%v]", c.GetRemoteAddr())
}

func (s *ServerHandler) OnReceive(c *socketx.SocketClient, data []byte, length int, from string) {
	log.Infof("web socket server received data [%s] length [%v] from [%v]", data, length, from)
	if string(data) == WEBSOCKET_DATA_PING {
		if _, err := c.Send([]byte(WEBSOCKET_DATA_PONG)); err != nil {
			log.Errorf(err.Error())
		}
	}
}

func (s *ServerHandler) OnClose(c *socketx.SocketClient) {
	log.Infof("connection [%v] closed", c.GetRemoteAddr())
}

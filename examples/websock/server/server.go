package main

import (
	"github.com/civet148/log"
	"github.com/civet148/socketx"
)

const (
	WEBSOCKET_DATA_PING = "ping"
	WEBSOCKET_DATA_PONG = "pong"
)

type Server struct {
	*socketx.SocketServer
}

func main() {
	var strUrl string
	//strUrl = "ws://0.0.0.0:6668/websocket" //web socket
	strUrl = "wss://0.0.0.0:6668/websocket?cert=cert.pem&key=key.pem" //web socket on TLS
	server(strUrl)
	var c = make(chan bool, 1)
	<-c //block main go routine
}

func server(strUrl string) {
	var server Server
	server.SocketServer = socketx.NewServer(strUrl)
	_ = server.Listen(&server)
}

func (s *Server) OnAccept(c *socketx.SocketClient) {
	log.Infof("connection accepted [%v]", c.GetRemoteAddr())
}

func (s *Server) OnReceive(c *socketx.SocketClient, data []byte, length int, from string) {
	log.Infof("web socket server received data [%s] length [%v] from [%v]", data, length, from)
	if string(data) == WEBSOCKET_DATA_PING {
		if _, err := c.Send([]byte(WEBSOCKET_DATA_PONG)); err != nil {
			log.Errorf(err.Error())
		}
	}
}

func (s *Server) OnClose(c *socketx.SocketClient) {
	log.Infof("connection [%v] closed", c.GetRemoteAddr())
}

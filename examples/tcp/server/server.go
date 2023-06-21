package main

import (
	"github.com/civet148/log"
	"github.com/civet148/socketx"
)

const (
	TCP_DATA_PING = "ping"
	TCP_DATA_PONG = "pong"
)

type Server struct {
	service *socketx.SocketServer
}

func main() {

	var url = "tcp://0.0.0.0:6666"
	server(url)

	var c = make(chan bool, 1)
	<-c //block main go routine
}

func server(strUrl string) {

	var server Server
	server.service = socketx.NewServer(strUrl)
	if err := server.service.Listen(&server); err != nil {
		log.Errorf(err.Error())
		return
	}
}

func (s *Server) OnAccept(c *socketx.SocketClient) {
	log.Infof("connection accepted [%v]", c.GetRemoteAddr())
}

func (s *Server) OnReceive(c *socketx.SocketClient, data []byte, length int, from string) {
	log.Infof("tcp server received data [%s] length [%v] from [%v]", data, length, from)
	if string(data) == TCP_DATA_PING {
		if _, err := c.Send([]byte(TCP_DATA_PONG)); err != nil {
			log.Errorf(err.Error())
		}
	}
}

func (s *Server) OnClose(c *socketx.SocketClient) {
	log.Infof("connection [%v] closed", c.GetRemoteAddr())
}

package main

import (
	"github.com/civet148/log"
	"github.com/civet148/socketx"
)

const (
	UDP_DATA_PING  = "ping"
	UDP_DATA_PONG  = "pong"
	UDP_SERVER_URL = "udp://0.0.0.0:6665"
)

type Server struct {
	service *socketx.SocketServer
}

func main() {

	server(UDP_SERVER_URL)

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
}

func (s *Server) OnReceive(c *socketx.SocketClient, data []byte, length int, from string) {
	log.Infof("udp server received data [%s] length [%v] from [%v] ", data, length, from)
	if string(data) == UDP_DATA_PING {
		if _, err := c.Send([]byte(UDP_DATA_PONG), from); err != nil {
			log.Errorf(err.Error())
		}
	}
}

func (s *Server) OnClose(c *socketx.SocketClient) {
}

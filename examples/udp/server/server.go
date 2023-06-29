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

type ServerHandler struct {
}

func init() {
	log.SetLevel("debug")
}

func main() {
	var handler ServerHandler
	sock := socketx.NewServer(UDP_SERVER_URL)
	if err := sock.Listen(&handler); err != nil {
		log.Errorf(err.Error())
		return
	}
}

func (s *ServerHandler) OnAccept(c *socketx.SocketClient) {
}

func (s *ServerHandler) OnReceive(c *socketx.SocketClient, data []byte, length int, from string) {
	log.Infof("udp server received data [%s] length [%v] from [%v] ", data, length, from)
	if string(data) == UDP_DATA_PING {
		if _, err := c.Send([]byte(UDP_DATA_PONG), from); err != nil {
			log.Errorf(err.Error())
		}
	}
}

func (s *ServerHandler) OnClose(c *socketx.SocketClient) {
}

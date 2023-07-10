package main

import (
	"github.com/civet148/log"
	"github.com/civet148/socketx"
	"github.com/civet148/socketx/api"
)

const (
	TCP_DATA_PING = "ping"
	TCP_DATA_PONG = "pong"
)

type ServerHandler struct {
}

func init() {
	log.SetLevel("debug")
}

func main() {

	var strUrl = "tcp://0.0.0.0:6666"
	var handler ServerHandler
	sock := socketx.NewServer(strUrl)
	if err := sock.Listen(&handler); err != nil {
		log.Errorf(err.Error())
		return
	}
}

func (s *ServerHandler) OnAccept(c *socketx.SocketClient) {
	log.Infof("connection accepted [%v]", c.GetRemoteAddr())
}

func (s *ServerHandler) OnReceive(c *socketx.SocketClient, msg *api.SockMessage) {
	data := msg.Data
	from := msg.From
	length := len(data)
	log.Infof("server received data [%s] length [%v] from [%v]", data, length, from)
	if string(data) == TCP_DATA_PING {
		if _, err := c.Send([]byte(TCP_DATA_PONG)); err != nil {
			log.Errorf(err.Error())
		}
	}
}

func (s *ServerHandler) OnClose(c *socketx.SocketClient) {
	log.Infof("connection [%v] closed", c.GetRemoteAddr())
}

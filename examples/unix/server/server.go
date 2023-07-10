package main

import (
	"github.com/civet148/log"
	"github.com/civet148/socketx"
	"github.com/civet148/socketx/api"
)

const (
	UNIX_SOCKET_URL = "unix:///tmp/unix.sock"
)
const (
	UNIX_DATA_PING = "ping"
	UNIX_DATA_PONG = "pong"
)

type ServerHandler struct {
}

func init() {
	log.SetLevel("debug")
}

func main() {
	var handler ServerHandler
	sock := socketx.NewServer(UNIX_SOCKET_URL)
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
	log.Infof("server received data [%s] length [%v] from [%v] ", data, length, from)
	if string(data) == UNIX_DATA_PING {
		if _, err := c.Send([]byte(UNIX_DATA_PONG)); err != nil {
			log.Errorf(err.Error())
		}
	}
}

func (s *ServerHandler) OnClose(c *socketx.SocketClient) {
	log.Infof("connection [%v] closed", c.GetRemoteAddr())
}

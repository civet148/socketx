package main

import (
	"github.com/civet148/log"
	"github.com/civet148/socketx"
	"time"
)

const (
	UNIX_SOCKET_URL = "unix:///tmp/unix.sock"
)

const (
	UNIX_DATA_PING = "ping"
	UNIX_DATA_PONG = "pong"
)

func init() {
	log.SetLevel("debug")
}

func main() {
	c := socketx.NewClient()
	if err := c.Connect(UNIX_SOCKET_URL); err != nil {
		log.Errorf(err.Error())
		return
	}

	for {
		if _, err := c.Send([]byte(UNIX_DATA_PING)); err != nil {
			log.Errorf(err.Error())
			break
		}

		msg, err := c.Recv(len(UNIX_DATA_PONG))
		if err != nil {
			log.Error(err.Error())
			break
		}
		data := msg.Data
		from := msg.From
		log.Infof("client received data [%s] length [%v] from [%v]", string(msg.Data), len(data), from)
		time.Sleep(3 * time.Second)
	}
}

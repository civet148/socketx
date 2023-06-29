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

		if data, from, err := c.Recv(len(UNIX_DATA_PONG)); err != nil {
			log.Error(err.Error())
			break
		} else {
			log.Infof("unix client received data [%s] length [%v] from [%v]", string(data), len(data), from)
		}

		time.Sleep(3 * time.Second)
	}
}

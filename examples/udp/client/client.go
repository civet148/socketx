package main

import (
	"github.com/civet148/log"
	"github.com/civet148/socketx"
	"time"
)

const (
	UDP_CLIENT_ADDR = "udp://127.0.0.1:6664"
	UDP_SERVER_ADDR = "udp://127.0.0.1:6665"
)

const (
	UDP_DATA_PING = "ping"
	UDP_DATA_PONG = "pong"
)

func init() {
	log.SetLevel("debug")
}

func main() {
	c := socketx.NewClient()
	if err := c.Listen(UDP_CLIENT_ADDR); err != nil {
		log.Errorf(err.Error())
		return
	}
	for {
		if _, err := c.Send([]byte(UDP_DATA_PING), UDP_SERVER_ADDR); err != nil {
			log.Errorf(err.Error())
			break
		}

		msg, err := c.Recv(-1)
		if err != nil {
			log.Error(err.Error())
			break
		}
		data := msg.Data
		from := msg.From
		log.Infof("client received data [%s] length [%v] from [%v]", string(msg.Data), len(data), from)
		time.Sleep(3 * time.Second)
	}
	return
}

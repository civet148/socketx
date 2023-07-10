package main

import (
	"github.com/civet148/log"
	"github.com/civet148/socketx"
	"time"
)

const (
	TCP_DATA_PING = "ping"
	TCP_DATA_PONG = "pong"
)

func init() {
	log.SetLevel("debug")
}

func main() {
	var strUrl = "tcp://127.0.0.1:6666"
	c := socketx.NewClient()
	if err := c.Connect(strUrl); err != nil {
		log.Errorf(err.Error())
		return
	}
	defer c.Close()

	for {
		if _, err := c.Send([]byte(TCP_DATA_PING)); err != nil {
			log.Errorf(err.Error())
			break
		}

		msg, err := c.Recv(len(TCP_DATA_PONG))
		if err != nil {
			log.Error(err.Error())
			break
		}
		data := msg.Data
		from := msg.From
		log.Infof("client received data [%s] length [%v] from [%v]", string(msg.Data), len(data), from)
		time.Sleep(1 * time.Second)
	}
}

package main

import (
	"github.com/civet148/log"
	"github.com/civet148/socketx"
	"github.com/civet148/socketx/api"
	"time"
)

const (
	WEBSOCKET_SERVER_URL = "wss://127.0.0.1:6668/websocket"
)

const (
	WEBSOCKET_DATA_PING = "ping"
	WEBSOCKET_DATA_PONG = "pong"
)

func init() {
	log.SetLevel("debug")
}

func main() {
	var err error
	c := socketx.NewClient()
	if err = c.Connect(WEBSOCKET_SERVER_URL); err != nil {
		log.Errorf(err.Error())
		return
	}

	for {
		if _, err = c.Send([]byte(WEBSOCKET_DATA_PING)); err != nil {
			log.Errorf(err.Error())
			break
		}
		var msg *api.SockMessage
		msg, err = c.Recv(-1)
		if err != nil {
			log.Errorf(err.Error())
			break
		}
		data := msg.Data
		from := msg.From
		log.Infof("client received data [%s] length [%v] from [%v]", string(msg.Data), len(data), from)
		time.Sleep(3 * time.Second)
	}
	return
}

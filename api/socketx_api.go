package api

import (
	"fmt"
	"github.com/civet148/gotools/parser"
	"github.com/civet148/log"
	"github.com/civet148/socketx/types"
)

type SockMessage struct {
	Sock    Socket //socket handle
	Data    []byte //data received
	From    string //remote address for UDP
	MsgType int    //only for websocket
}

type Socket interface {
	Listen() (err error)                               // bind and listen on address and port
	Accept() Socket                                    // accept connection...
	Connect() (err error)                              // for tcp/web socket
	Send(data []byte, to ...string) (n int, err error) // send to...
	Recv(length int) (msg *SockMessage, err error)     // receive from... if length > 0, will receive the bytes specified.
	Close() (err error)                                // close socket
	GetLocalAddr() string                              // get socket local address
	GetRemoteAddr() string                             // get socket remote address
	GetSocketType() types.SocketType                   // get socket type
}

type SocketInstance func(ui *parser.UrlInfo) Socket

var instances = make(map[types.SocketType]SocketInstance)

func Register(sockType types.SocketType, inst SocketInstance) (err error) {
	if _, ok := instances[sockType]; !ok {

		instances[sockType] = inst
		return
	}
	err = fmt.Errorf("socket type [%v] instance already exists", sockType)
	log.Errorf("%v", err.Error())
	return
}

func NewSocketInstance(sockType types.SocketType, ui *parser.UrlInfo) (s Socket) {
	if inst, ok := instances[sockType]; !ok {
		log.Errorf("socket type [%v] instance not register", sockType)
		return nil
	} else {
		s = inst(ui)
	}
	return
}

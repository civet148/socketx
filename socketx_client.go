package socketx

import (
	"fmt"
	"github.com/civet148/log"
	"github.com/civet148/socketx/api"
	_ "github.com/civet148/socketx/tcpsock"  //register TCP instance
	_ "github.com/civet148/socketx/udpsock"  //register UDP instance
	_ "github.com/civet148/socketx/unixsock" //register UNIX instance
	_ "github.com/civet148/socketx/websock"  //register WEBSOCKET instance
)

type SocketClient struct {
	sock   api.Socket
	closed bool
}

func init() {
	log.SetLevel("info")
}

func NewClient() *SocketClient {
	return &SocketClient{}
}

// IPv4      => 		tcp://127.0.0.1:6666 [tcp4://127.0.0.1:6666]
// WebSocket => 		ws://127.0.0.1:6668 [wss://127.0.0.1:6668]
func (w *SocketClient) Connect(url string) (err error) {
	var s api.Socket
	if s = CreateSocket(url); s == nil {
		return fmt.Errorf("create socket by url [%v] failed", url)
	}
	w.sock = s
	return w.sock.Connect()
}

// only for UDP
func (w *SocketClient) Listen(url string) (err error) {
	if w.sock = CreateSocket(url); w.sock == nil {
		return fmt.Errorf("create socket by url [%v] failed", url)
	}
	return w.sock.Listen()
}

func (w *SocketClient) Send(data []byte, to ...string) (n int, err error) {
	return w.send(w.sock, data, to...)
}

func (w *SocketClient) Recv(length int) (data []byte, from string, err error) {
	return w.recv(w.sock, length)
}

func (w *SocketClient) GetLocalAddr() (addr string) {
	return w.sock.GetLocalAddr()
}

func (w *SocketClient) GetRemoteAddr() (addr string) {
	return w.sock.GetRemoteAddr()
}

func (w *SocketClient) Close() (err error) {
	return w.sock.Close()
}

func (w *SocketClient) IsClosed() bool {
	return w.closed
}

func (w *SocketClient) send(s api.Socket, data []byte, to ...string) (n int, err error) {
	return s.Send(data, to...)
}

func (w *SocketClient) recv(s api.Socket, length int) (data []byte, from string, err error) {
	return s.Recv(length)
}

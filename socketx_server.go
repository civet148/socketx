package socketx

import (
	"fmt"
	"github.com/civet148/gotools/parser"
	"github.com/civet148/log"
	"github.com/civet148/socketx/api"
	"github.com/civet148/socketx/types"
	"strings"
	"sync"
)

type SocketHandler interface {
	OnAccept(c *SocketClient)
	OnReceive(c *SocketClient, data []byte, length int, from string)
	OnClose(c *SocketClient)
}

type SocketServer struct {
	url       string                       //listen url
	sock      api.Socket                   //server socket
	handler   SocketHandler                //server callback handler
	accepting chan api.Socket              //client connection accepted
	receiving chan api.Socket              //client message received
	quiting   chan api.Socket              //client connection closed
	clients   map[api.Socket]*SocketClient //socket clients
	locker    *sync.Mutex                  //locker mutex
	done      chan bool                    //force close server socket
}

func init() {
	log.SetLevel("info")
}

func NewServer(url string) *SocketServer {

	var s api.Socket
	s = CreateSocket(url)

	return &SocketServer{
		url:       url,
		locker:    &sync.Mutex{},
		sock:      s,
		done:      make(chan bool),
		accepting: make(chan api.Socket, 1000),
		quiting:   make(chan api.Socket, 1000),
		clients:   make(map[api.Socket]*SocketClient, 0),
	}
}

// TCP       => 		tcp://127.0.0.1:6666
// UDP       => 		udp://127.0.0.1:6667
// WebSocket => 		ws://127.0.0.1:6668/ wss://127.0.0.1:6668/websocket?cert=cert.pem&key=key.pem
func (w *SocketServer) Listen(handler SocketHandler) (err error) {
	w.handler = handler
	if err = w.sock.Listen(); err != nil {
		log.Errorf(err.Error())
		return
	}
	log.Infof("listen [%v] address [%v] ok", w.sock.GetSocketType(), w.sock.GetLocalAddr())
	if w.sock.GetSocketType() != types.SocketType_UDP {
		go func() {
			//log.Tracef("start goroutine for channel event accepting/quiting")
			for {
				select {
				case s := <-w.accepting: //client connection coming...
					w.onAccept(s)
				case s := <-w.quiting: //client connection closed
					w.onClose(s)
					//default: //disable default because of high CPU performance
				}
			}
		}()

		//new go routine for accept new connections
		go func() {
			//log.Tracef("start goroutine for accept new connection")
			for {
				if s := w.sock.Accept(); s != nil { //socket accepting...
					w.accepting <- s
				}
			}
		}()
	} else {
		w.onAccept(w.sock)
	}

	<-w.done //wait for signal
	return
}

func (w *SocketServer) Close() {
	w.done <- true
	w.sock.Close()
	w.closeClientAll()
}

func (w *SocketServer) CloseClient(client *SocketClient) (err error) {
	return w.closeSocket(client.sock)
}

func (w *SocketServer) Send(client *SocketClient, data []byte, to ...string) (n int, err error) {
	return w.sendSocket(client.sock, data, to...)
}

func (w *SocketServer) GetClientCount() int {
	return w.getClientCount()
}

func (w *SocketServer) GetClientAll() (clients []*SocketClient) {
	return w.getClientAll()
}

func (w *SocketServer) closeSocket(s api.Socket) (err error) {
	if s == nil {
		return fmt.Errorf("close socket is nil")
	}
	w.removeClient(s)
	return s.Close()
}

func (w *SocketServer) sendSocket(s api.Socket, data []byte, to ...string) (n int, err error) {
	if s == nil || len(data) == 0 {
		err = fmt.Errorf("send socket is nil or data length is 0")
		return
	}
	return s.Send(data, to...)
}

func (w *SocketServer) recvSocket(s api.Socket) (data []byte, from string, err error) {
	if s == nil {
		err = fmt.Errorf("send socket is nil")
		return
	}
	return s.Recv(-1)
}

func (w *SocketServer) onAccept(s api.Socket) {
	c := w.addClient(s)
	go w.readSocket(s)
	w.handler.OnAccept(c)
}

func (w *SocketServer) onClose(s api.Socket) {
	w.handler.OnClose(w.removeClient(s))
}

func (w *SocketServer) onReceive(s api.Socket, data []byte, length int, from string) {
	c := w.getClient(s)
	w.handler.OnReceive(c, data, length, from)
}

func (w *SocketServer) readSocket(s api.Socket) {
	for {
		if data, from, err := w.recvSocket(s); err != nil {
			w.quiting <- s
			break
		} else if len(data) > 0 {
			w.onReceive(s, data, len(data), from)
		}
	}
}

func (w *SocketServer) lock() {
	w.locker.Lock()
}

func (w *SocketServer) unlock() {
	w.locker.Unlock()
}

func (w *SocketServer) closeClientAll() {
	w.lock()
	defer w.unlock()
	for s, _ := range w.clients {
		_ = s.Close()
		delete(w.clients, s)
	}
}

func (w *SocketServer) addClient(s api.Socket) (client *SocketClient) {
	client = &SocketClient{
		sock: s,
	}
	w.lock()
	defer w.unlock()
	w.clients[client.sock] = client
	return client
}

func (w *SocketServer) removeClient(s api.Socket) (client *SocketClient) {
	w.lock()
	defer w.unlock()
	client = w.clients[s]
	delete(w.clients, s)
	return
}

func (w *SocketServer) getClient(s api.Socket) (client *SocketClient) {
	var ok bool
	w.lock()
	defer w.unlock()
	if client, ok = w.clients[s]; ok {
		return
	}
	return
}

func (w *SocketServer) getClientCount() int {
	w.lock()
	defer w.unlock()
	return len(w.clients)
}

func (w *SocketServer) getClientAll() (clients []*SocketClient) {
	w.lock()
	defer w.unlock()
	for _, v := range w.clients {
		clients = append(clients, v)
	}
	return
}

func CreateSocket(url string) (s api.Socket) {

	url = strings.ToLower(url)
	ui := parser.ParseUrl(url)
	switch ui.Scheme {
	case types.URL_SCHEME_TCP, types.URL_SCHEME_TCP4, types.URL_SCHEME_TCP6:
		s = api.NewSocketInstance(types.SocketType_TCP, ui)
	case types.URL_SCHEME_WS, types.URL_SCHEME_WSS:
		s = api.NewSocketInstance(types.SocketType_WEB, ui)
	case types.URL_SCHEME_UDP, types.URL_SCHEME_UDP4, types.URL_SCHEME_UDP6:
		s = api.NewSocketInstance(types.SocketType_UDP, ui)
	case types.URL_SCHEME_UNIX:
		s = api.NewSocketInstance(types.SocketType_UNIX, ui)
	default:
		{
			url = types.URL_SCHEME_TCP + parser.URL_SCHEME_SEP + url
			ui = parser.ParseUrl(url)
			s = api.NewSocketInstance(types.SocketType_TCP, ui) //default 'tcp'
		}
	}
	return
}

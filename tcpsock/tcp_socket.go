package tcpsock

import (
	"fmt"
	"github.com/civet148/gotools/parser"
	"github.com/civet148/log"
	"github.com/civet148/socketx/api"
	"github.com/civet148/socketx/types"
	"net"
	"sync"
)

type socket struct {
	ui       *parser.UrlInfo
	conn     net.Conn
	listener net.Listener
	closed   bool
	locker   sync.RWMutex
}

func init() {
	_ = api.Register(types.SocketType_TCP, NewSocket)
}

func NewSocket(ui *parser.UrlInfo) api.Socket {

	return &socket{
		ui: ui,
	}
}

func (s *socket) Listen() (err error) {
	var network = s.getNetwork()
	strAddr := s.ui.GetHost()
	//log.Debugf("trying listen [%v] protocol [%v]", strAddr, s.ui.GetScheme())
	s.listener, err = net.Listen(network, strAddr)
	if err != nil {
		log.Errorf("listen tcp address [%s] failed", strAddr)
		return
	}
	return
}

func (s *socket) Accept() api.Socket {
	conn, err := s.listener.Accept()
	if err != nil {
		return nil
	}
	return &socket{
		conn: conn,
	}
}

func (s *socket) Connect() (err error) {
	var network = s.getNetwork()
	addr := s.ui.GetHost()
	var tcpAddr *net.TCPAddr
	tcpAddr, err = net.ResolveTCPAddr(network, addr)
	if err != nil {
		log.Errorf("resolve tcp address [%s] failed, error [%s]", addr, err)
		return err
	}

	s.conn, err = net.DialTCP(network, nil, tcpAddr)
	if err != nil {
		log.Errorf("dial tcp to [%s] failed", addr)
		return err
	}
	return
}

func (s *socket) Send(data []byte, to ...string) (n int, err error) {
	s.locker.Lock()
	defer s.locker.Unlock()
	return s.conn.Write(data)
}

// length <= 0, default PACK_FRAGMENT_MAX=1500 bytes
func (s *socket) Recv(length int) (msg *api.SockMessage, err error) {

	var once bool
	var recv, left int
	if length <= 0 {
		once = true
		length = types.PACK_FRAGMENT_MAX
	}
	left = length
	data := s.makeBuffer(length)
	var n int
	if once {
		if n, err = s.conn.Read(data); err != nil {
			log.Errorf("read data error [%v]", err.Error())
			return
		}
		recv = n
	} else {

		for left > 0 {
			if n, err = s.conn.Read(data[recv:]); err != nil {
				log.Errorf("read data error [%v]", err.Error())
				return
			}
			left -= n
			recv += n
		}
	}

	if recv < length {
		data = data[:recv]
	}
	from := s.conn.RemoteAddr().String()
	return &api.SockMessage{
		Sock: s,
		Data: data,
		From: from,
	}, nil
}

func (s *socket) Close() (err error) {
	if s.closed {
		err = fmt.Errorf("socket already closed")
		return
	}
	if s.conn == nil {
		err = fmt.Errorf("socket is nil")
		log.Error(err.Error())
		return
	}
	s.closed = true
	return s.conn.Close()
}

func (s *socket) GetLocalAddr() string {
	if s.conn == nil {
		return s.ui.GetHost()
	}
	return s.conn.LocalAddr().String()
}

func (s *socket) GetRemoteAddr() string {
	if s.conn == nil {
		return ""
	}
	return s.conn.RemoteAddr().String()
}

func (s *socket) GetSocketType() types.SocketType {
	return types.SocketType_TCP
}

func (s *socket) getNetwork() string {
	if s.isTcp6() {
		return types.NETWORK_TCPv6
	}
	return types.NETWORK_TCPv4
}

func (s *socket) isTcp6() (ok bool) {
	scheme := s.ui.GetScheme()
	if scheme == types.URL_SCHEME_TCP6 {
		return true
	}
	return
}

func (s *socket) makeBuffer(length int) []byte {
	return make([]byte, length)
}

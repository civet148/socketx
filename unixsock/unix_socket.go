package unixsock

import (
	"fmt"
	"github.com/civet148/gotools/parser"
	"github.com/civet148/log"
	"github.com/civet148/socketx/api"
	"github.com/civet148/socketx/types"
	"net"
	"os"
	"strings"
	"sync"
)

type socket struct {
	ui       *parser.UrlInfo
	conn     net.Conn
	listener *net.UnixListener
	closed   bool
	locker   sync.RWMutex
}

func init() {
	_ = api.Register(types.SocketType_UNIX, NewSocket)
}

func NewSocket(ui *parser.UrlInfo) api.Socket {

	return &socket{
		ui: ui,
	}
}

func (s *socket) Listen() (err error) {
	var fi os.FileInfo
	var network = s.getNetwork()
	addr := s.getUnixSockFile()
	fi, err = os.Stat(addr)
	if err != nil && !os.IsNotExist(err) {
		return log.Errorf(err.Error())
	}
	if fi.Name() != "" {
		if err = os.Remove(addr); err != nil {
			return log.Errorf("remove file error [%v]", err.Error())
		}
	}

	var unixAddr *net.UnixAddr
	unixAddr, err = net.ResolveUnixAddr(network, s.ui.GetPath())
	if err != nil {
		return log.Errorf("resolve unix addr %s error [%s]", s.ui.GetPath(), err.Error())
	}
	if s.listener, err = net.ListenUnix("unix", unixAddr); err != nil {
		return log.Errorf("listen tcp address [%s] error [%s]", addr, err.Error())
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
		ui:   s.ui,
	}
}

func (s *socket) Connect() (err error) {
	var network = s.getNetwork()
	addr := s.getUnixSockFile()
	var unixAddr *net.UnixAddr
	unixAddr, err = net.ResolveUnixAddr(network, addr)
	if err != nil {
		return log.Errorf("resolve address [%s] error [%s]", addr, err)
	}

	s.conn, err = net.DialUnix(network, nil, unixAddr)
	if err != nil {
		return log.Errorf("dial [%s] to [%s] error [%s]", network, addr, err.Error())
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
			return nil, log.Errorf("read data from [%s] error [%v]", s.GetRemoteAddr(), err.Error())
		}
		recv = n
	} else {

		for left > 0 {
			if n, err = s.conn.Read(data[recv:]); err != nil {
				return nil, log.Errorf("read data from [%s] error [%v]", s.GetRemoteAddr(), err.Error())
			}
			left -= n
			recv += n
		}
	}

	if recv < length {
		data = data[:recv]
	}
	from := s.GetLocalAddr()
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
		return log.Error("socket is nil")
	}
	s.closed = true
	return s.conn.Close()
}

func (s *socket) GetLocalAddr() (strAddr string) {
	return s.getUnixSockFile()
}

func (s *socket) GetRemoteAddr() (strAddr string) {
	return s.getUnixSockFile()
}

func (s *socket) GetSocketType() types.SocketType {
	return types.SocketType_UNIX
}

func (s *socket) getUnixSockFile() (strSockFile string) {

	if s.ui == nil {
		return
	}
	strSockFile = s.ui.GetPath()
	if !strings.HasSuffix(strSockFile, "sock") {
		panic("unix socket must .sock as file suffix")
	}
	return
}

func (s *socket) getNetwork() string {
	return types.NETWORK_UNIX
}

func (s *socket) makeBuffer(length int) []byte {
	return make([]byte, length)
}

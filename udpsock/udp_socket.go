package udpsock

import (
	"fmt"
	"github.com/civet148/gotools/parser"
	"github.com/civet148/log"
	"github.com/civet148/socketx/api"
	"github.com/civet148/socketx/types"
	"net"
	"strings"
	"sync"
)

type socket struct {
	ui     *parser.UrlInfo
	conn   *net.UDPConn
	closed bool
	locker sync.RWMutex
}

func init() {
	_ = api.Register(types.SocketType_UDP, NewSocket)
}

func NewSocket(ui *parser.UrlInfo) api.Socket {

	return &socket{
		ui: ui,
	}
}

func (s *socket) Listen() (err error) {

	var strAddr = s.ui.GetHost()
	var udpAddr *net.UDPAddr
	var network = s.getNetwork()

	if udpAddr, err = net.ResolveUDPAddr(network, strAddr); err != nil {
		return log.Errorf("resolve UDP addr [%v] error [%v]", strAddr, err.Error())
	}

	if s.conn, err = net.ListenUDP(network, udpAddr); err != nil {
		return log.Errorf("listen UDP addr [%v] error [%v]", strAddr, err.Error())
	}
	return
}

func (s *socket) Accept() api.Socket {
	log.Warnf("accept method not for UDP socket")
	return nil
}

func (s *socket) Connect() (err error) {
	return fmt.Errorf("only for TCP/WEB socket")
}

func (s *socket) Send(data []byte, to ...string) (n int, err error) {

	var udpAddr *net.UDPAddr
	var network = s.getNetwork()

	if len(to) == 0 {
		return 0, fmt.Errorf("UDP send method to parameter required")
	}
	s.locker.Lock()
	defer s.locker.Unlock()
	strToAddr := to[0]
	nSep := len(parser.URL_SCHEME_SEP)
	if strings.Contains(strToAddr, parser.URL_SCHEME_SEP) {
		nIndex := strings.Index(strToAddr, parser.URL_SCHEME_SEP)
		strToAddr = strToAddr[nIndex+nSep:]
	}

	if udpAddr, err = net.ResolveUDPAddr(network, strToAddr); err != nil {
		return 0, log.Errorf("resolve UDP addr [%v] error [%v]", strToAddr, err.Error())
	}
	return s.conn.WriteToUDP(data, udpAddr)
}

func (s *socket) Recv(length int) (msg *api.SockMessage, err error) {
	var n int
	var udpAddr *net.UDPAddr
	data := s.makeBuffer(types.PACK_FRAGMENT_MAX)
	if n, udpAddr, err = s.conn.ReadFromUDP(data); err != nil {
		log.Errorf("read from UDP error [%v]", err.Error())
		return
	}
	return &api.SockMessage{
		Sock: s,
		Data: data[:n],
		From: udpAddr.String(),
	}, nil
}

func (s *socket) Close() (err error) {
	if s.closed {
		return fmt.Errorf("socket already closed")
	}
	if s.conn == nil {
		return fmt.Errorf("socket is nil")
	}
	s.closed = true
	return s.conn.Close()
}

func (s *socket) GetLocalAddr() string {
	return s.conn.LocalAddr().String()
}

func (s *socket) GetRemoteAddr() (addr string) {
	return
}

func (s *socket) GetSocketType() types.SocketType {
	return types.SocketType_UDP
}

func (s *socket) getNetwork() string {
	if s.isUDP6() {
		return types.NETWORK_UDPv6
	}
	return types.NETWORK_UDPv4
}

func (s *socket) isUDP6() (ok bool) {
	scheme := s.ui.GetScheme()
	if scheme == types.URL_SCHEME_UDP6 {
		return true
	}
	return
}

func (s *socket) makeBuffer(length int) []byte {
	return make([]byte, length)
}

package types

const (
	URL_SCHEME_TCP  = "tcp"
	URL_SCHEME_TCP4 = "tcp4"
	URL_SCHEME_TCP6 = "tcp6"
	URL_SCHEME_UDP  = "udp"
	URL_SCHEME_UDP4 = "udp4"
	URL_SCHEME_UDP6 = "udp6"
	URL_SCHEME_WS   = "ws"
	URL_SCHEME_WSS  = "wss"
	URL_SCHEME_UNIX = "unix"
)

const (
	PACK_FRAGMENT_MAX = 1500
	TCP_FRAGMENT_MAX  = 1024 * 1024
)

const (
	NETWORK_TCP   = "tcp"
	NETWORK_TCPv4 = "tcp4"
	NETWORK_TCPv6 = "tcp6"
	NETWORK_UDP   = "udp"
	NETWORK_UDPv4 = "udp4"
	NETWORK_UDPv6 = "udp6"
	NETWORK_UNIX  = "unix"
)

const (
	WSS_TLS_CERT = "cert"
	WSS_TLS_KEY  = "key"
)

type SocketType int

const (
	SocketType_TCP  SocketType = 1
	SocketType_WEB  SocketType = 2
	SocketType_UDP  SocketType = 3
	SocketType_UNIX SocketType = 4
)

func (s SocketType) GoString() string {
	return s.String()
}

func (s SocketType) String() string {
	switch s {
	case SocketType_TCP:
		return "TCP"
	case SocketType_WEB:
		return "WEBSOCKET"
	case SocketType_UDP:
		return "UDP"
	case SocketType_UNIX:
		return "UNIX"
	}
	return "SocketType<Unknown>"
}

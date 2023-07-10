# a socket wrapper for TCP/UDP/WEB/UNIX socket

# 1. TCP socket

## 1.1 TCP client 

```go
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
```
## 1.2 TCP server 

```go
package main

import (
	"github.com/civet148/log"
	"github.com/civet148/socketx"
	"github.com/civet148/socketx/api"
)

const (
	TCP_DATA_PING = "ping"
	TCP_DATA_PONG = "pong"
)

type ServerHandler struct {
}

func init() {
	log.SetLevel("debug")
}

func main() {

	var strUrl = "tcp://0.0.0.0:6666"
	var handler ServerHandler
	sock := socketx.NewServer(strUrl)
	if err := sock.Listen(&handler); err != nil {
		log.Errorf(err.Error())
		return
	}
}

func (s *ServerHandler) OnAccept(c *socketx.SocketClient) {
	log.Infof("connection accepted [%v]", c.GetRemoteAddr())
}

func (s *ServerHandler) OnReceive(c *socketx.SocketClient, msg *api.SockMessage) {
	data := msg.Data
	from := msg.From
	length := len(data)
	log.Infof("server received data [%s] length [%v] from [%v]", data, length, from)
	if string(data) == TCP_DATA_PING {
		if _, err := c.Send([]byte(TCP_DATA_PONG)); err != nil {
			log.Errorf(err.Error())
		}
	}
}

func (s *ServerHandler) OnClose(c *socketx.SocketClient) {
	log.Infof("connection [%v] closed", c.GetRemoteAddr())
}
```

# 2. UDP socket

## 2.1 UDP client

```go
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
```
## 2.2 UDP server 

```go
package main

import (
	"github.com/civet148/log"
	"github.com/civet148/socketx"
	"github.com/civet148/socketx/api"
)

const (
	UDP_DATA_PING  = "ping"
	UDP_DATA_PONG  = "pong"
	UDP_SERVER_URL = "udp://0.0.0.0:6665"
)

type ServerHandler struct {
}

func init() {
	log.SetLevel("debug")
}

func main() {
	var handler ServerHandler
	sock := socketx.NewServer(UDP_SERVER_URL)
	if err := sock.Listen(&handler); err != nil {
		log.Errorf(err.Error())
		return
	}
}

func (s *ServerHandler) OnAccept(c *socketx.SocketClient) {
}

func (s *ServerHandler) OnReceive(c *socketx.SocketClient, msg *api.SockMessage) {
	data := msg.Data
	from := msg.From
	length := len(data)
	log.Infof("server received data [%s] length [%v] from [%v] ", data, length, from)
	if string(data) == UDP_DATA_PING {
		if _, err := c.Send([]byte(UDP_DATA_PONG), from); err != nil {
			log.Errorf(err.Error())
		}
	}
}

func (s *ServerHandler) OnClose(c *socketx.SocketClient) {
}
```

# 3. WEB socket 

## 3.1 WebSocket client 

```go
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

```
## 3.2 WebSocket server 

```go
package main

import (
	"github.com/civet148/log"
	"github.com/civet148/socketx"
	"github.com/civet148/socketx/api"
)

const (
	//WEBSOCKET_SERVER_URL = "ws://0.0.0.0:6668/websocket"
	WEBSOCKET_SERVER_URL = "wss://0.0.0.0:6668/websocket?cert=cert.pem&key=key.pem"
)

const (
	WEBSOCKET_DATA_PING = "ping"
	WEBSOCKET_DATA_PONG = "pong"
)

type ServerHandler struct {
}

func main() {
	log.SetLevel("debug")
	var handler ServerHandler
	sock := socketx.NewServer(WEBSOCKET_SERVER_URL)
	_ = sock.Listen(&handler)
}

func (s *ServerHandler) OnAccept(c *socketx.SocketClient) {
	log.Infof("connection accepted [%v]", c.GetRemoteAddr())
}

func (s *ServerHandler) OnReceive(c *socketx.SocketClient, msg *api.SockMessage) {
	data := msg.Data
	from := msg.From
	length := len(data)
	log.Infof("server received data [%s] length [%v] from [%v] type [%v]", data, length, from, msg.MsgType)
	if string(data) == WEBSOCKET_DATA_PING {
		if _, err := c.Send([]byte(WEBSOCKET_DATA_PONG)); err != nil {
			log.Errorf(err.Error())
		}
	}
}

func (s *ServerHandler) OnClose(c *socketx.SocketClient) {
	log.Infof("connection [%v] closed", c.GetRemoteAddr())
}

```

# 4. UNIX socket

## 4.1 UNIX client 

```go
package main

import (
	"github.com/civet148/log"
	"github.com/civet148/socketx"
	"time"
)

const (
	UNIX_SOCKET_URL = "unix:///tmp/unix.sock"
)

const (
	UNIX_DATA_PING = "ping"
	UNIX_DATA_PONG = "pong"
)

func init() {
	log.SetLevel("debug")
}

func main() {
	c := socketx.NewClient()
	if err := c.Connect(UNIX_SOCKET_URL); err != nil {
		log.Errorf(err.Error())
		return
	}

	for {
		if _, err := c.Send([]byte(UNIX_DATA_PING)); err != nil {
			log.Errorf(err.Error())
			break
		}

		msg, err := c.Recv(len(UNIX_DATA_PONG))
		if err != nil {
			log.Error(err.Error())
			break
		}
		data := msg.Data
		from := msg.From
		log.Infof("client received data [%s] length [%v] from [%v]", string(msg.Data), len(data), from)
		time.Sleep(3 * time.Second)
	}
}

```
## 4.2 UNIX server 

```go
package main

import (
	"github.com/civet148/log"
	"github.com/civet148/socketx"
	"github.com/civet148/socketx/api"
)

const (
	UNIX_SOCKET_URL = "unix:///tmp/unix.sock"
)
const (
	UNIX_DATA_PING = "ping"
	UNIX_DATA_PONG = "pong"
)

type ServerHandler struct {
}

func init() {
	log.SetLevel("debug")
}

func main() {
	var handler ServerHandler
	sock := socketx.NewServer(UNIX_SOCKET_URL)
	if err := sock.Listen(&handler); err != nil {
		log.Errorf(err.Error())
		return
	}
}

func (s *ServerHandler) OnAccept(c *socketx.SocketClient) {
	log.Infof("connection accepted [%v]", c.GetRemoteAddr())
}

func (s *ServerHandler) OnReceive(c *socketx.SocketClient, msg *api.SockMessage) {
	data := msg.Data
	from := msg.From
	length := len(data)
	log.Infof("server received data [%s] length [%v] from [%v] ", data, length, from)
	if string(data) == UNIX_DATA_PING {
		if _, err := c.Send([]byte(UNIX_DATA_PONG)); err != nil {
			log.Errorf(err.Error())
		}
	}
}

func (s *ServerHandler) OnClose(c *socketx.SocketClient) {
	log.Infof("connection [%v] closed", c.GetRemoteAddr())
}

```


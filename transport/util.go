package transport

import (
	"net"
)



type TargetConn struct {
	net.Conn
	TargetAddr string
}

type Transport interface {
	Dial(addr, targetAddr string) (net.Conn, error)
	Listen(addr string) (<-chan *TargetConn, error)
}





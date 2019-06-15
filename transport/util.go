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

func NewTransport(network string) Transport {
	switch network {
	case "tcp":
		return NewTCP()
	case "socks5", "socks5h":
		return NewSocks5()
	case "quic":
		return NewQUIC()
	case "kcp":
		return NewKCP()
	default:
		panic("invalid transport type: " + network)
	}
}

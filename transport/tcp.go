package transport

import (
	"net"
	"time"

	"github.com/golang/glog"
	"github.com/wweir/sower/transport/parser"
)

type tcp struct {
	DialTimeout time.Duration
}

func NewTCP() Transport {
	return &tcp{DialTimeout: 5 * time.Second}
}

func (t *tcp) Dial(addr, targetAddr string) (net.Conn, error) {
	conn, err := net.DialTimeout("tcp", addr, t.DialTimeout)
	if err != nil {
		return nil, err
	}

	conn.(*net.TCPConn).SetKeepAlive(true)

	return parser.WithTarget(conn, targetAddr)
}

func (t *tcp) Listen(addr string) (<-chan *TargetConn, error) {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	connCh := make(chan *TargetConn)
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				glog.Fatalln("TCP listen:", err)
			}
			conn.(*net.TCPConn).SetKeepAlive(true)

			c, addr, err := parser.ParseAddr(conn)
			if err != nil {
				glog.Errorln("parse addr:", err)
			}
			connCh <- &TargetConn{c, addr}
		}
	}()
	return connCh, nil
}

package transport

import (
	"net"
	"time"

	"github.com/golang/glog"
	"github.com/wweir/sower/transport/router"
)

type tcp struct {
	raddr string

	DialTimeout time.Duration
}

func NewTCP(raddr string) Transport {
	return &tcp{raddr: raddr, DialTimeout: 5 * time.Second}
}

func (t *tcp) Dial(addr string) (net.Conn, error) {
	conn, err := net.DialTimeout("tcp", t.raddr, t.DialTimeout)
	if err != nil {
		return nil, err
	}

	conn.(*net.TCPConn).SetKeepAlive(true)

	return router.WriteAddr(conn, addr)
}

func (t *tcp) Listen(addr string) (<-chan *router.TargetConn, error) {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	connCh := make(chan *router.TargetConn)
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				glog.Fatalln("TCP listen:", err)
			}
			conn.(*net.TCPConn).SetKeepAlive(true)

			if tgtConn, err := router.ParseAddr(conn); err != nil {
				glog.Errorln("parse addr:", err)
			} else {
				connCh <- tgtConn
			}
		}
	}()
	return connCh, nil
}

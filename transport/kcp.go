package transport

import (
	"net"

	"github.com/golang/glog"
	"github.com/pkg/errors"
	"github.com/wweir/sower/transport/router"
	kcp "github.com/xtaci/kcp-go"
)

type kcpTran struct {
	client
	server
}
type client struct {
	raddr string

	DataShard    int
	ParityShard  int
	DSCP         int
	SockBuf      int
	AckNodelay   bool
	NoDelay      int
	Interval     int
	Resend       int
	NoCongestion int
	SndWnd       int
	RcvWnd       int
	MTU          int
}
type server struct {
	DataShard   int
	ParityShard int
	DSCP        int
	SockBuf     int
}

func NewKCP(raddr string) Transport {
	return &kcpTran{
		client: client{
			raddr:        raddr,
			DataShard:    10,
			ParityShard:  3,
			DSCP:         0,
			SockBuf:      4194304,
			NoDelay:      0,
			Interval:     50,
			Resend:       0,
			NoCongestion: 0,
			SndWnd:       0,
			RcvWnd:       0,
			MTU:          1350,
		},
		server: server{
			DataShard:   10,
			ParityShard: 3,
			DSCP:        0,
			SockBuf:     4194304,
		},
	}
}

func (c *client) Dial(addr string) (net.Conn, error) {
	conn, err := kcp.DialWithOptions(c.raddr, nil, c.DataShard, c.ParityShard)
	if err != nil {
		return nil, errors.Wrap(err, "dial")
	}

	conn.SetStreamMode(true)
	conn.SetWriteDelay(false)
	conn.SetNoDelay(c.NoDelay, c.Interval, c.Resend, c.NoCongestion)
	conn.SetWindowSize(c.SndWnd, c.RcvWnd)
	conn.SetMtu(c.MTU)
	conn.SetACKNoDelay(c.AckNodelay)

	if err := conn.SetDSCP(c.DSCP); err != nil {
		return nil, errors.Wrap(err, "SetDSCP")
	}
	if err := conn.SetReadBuffer(c.SockBuf); err != nil {
		return nil, errors.Wrap(err, "SetReadBuffer")
	}
	if err := conn.SetWriteBuffer(c.SockBuf); err != nil {
		return nil, errors.Wrap(err, "SetWriteBuffer")
	}

	return router.WriteAddr(conn, addr)
}

func (s *server) Listen(port string) (<-chan *router.TargetConn, error) {
	ln, err := kcp.ListenWithOptions(port, nil, s.DataShard, s.ParityShard)
	if err != nil {
		return nil, err
	}

	if err := ln.SetDSCP(s.DSCP); err != nil {
		return nil, errors.Wrap(err, "SetDSCP")
	}
	if err := ln.SetReadBuffer(s.SockBuf); err != nil {
		return nil, errors.Wrap(err, "SetReadBuffer")
	}
	if err := ln.SetWriteBuffer(s.SockBuf); err != nil {
		return nil, errors.Wrap(err, "SetWriteBuffer")
	}

	connCh := make(chan *router.TargetConn)
	go func() {
		for {
			conn, err := ln.AcceptKCP()
			if err != nil {
				glog.Fatalln("KCP listen:", err)
			}

			if tgtConn, err := router.ParseAddr(conn); err != nil {
				glog.Errorln("parse addr:", err)
			} else {
				connCh <- tgtConn
			}
		}
	}()

	return connCh, nil
}

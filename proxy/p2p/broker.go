package p2p

import (
	"encoding/binary"
	"net"
	"sync"

	"github.com/golang/glog"
)

type broker struct {
	readHeader  header
	writeHeader header
	cache       sync.Map
}

func (b *broker) Serve(port string) error {
	ln, err := net.Listen("tcp", port)
	if err != nil {
		return err
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			return err
		}

		go b.Handle(conn)
	}
}

func (b *broker) Handle(conn net.Conn) {
	if err := binary.Read(conn, binary.BigEndian, b.readHeader); err != nil {
		glog.Errorln("read init header", err)
		return
	}

	switch b.readHeader.command {
	case p2s:
		data := new(p2sT)
		if err := binary.Read(conn, binary.BigEndian, data); err != nil {
			glog.Errorln("read init header", err)
			return
		}
		b.cache.Store(data.uuid.String(), conn)

	case a2s:
		data := new(p2sT)
		if err := binary.Read(conn, binary.BigEndian, data); err != nil {
			glog.Errorln("read init header", err)
			return
		}

	case b2s:
	}
}

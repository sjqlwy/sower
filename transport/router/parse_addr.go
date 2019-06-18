package router

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"

	"github.com/wweir/sower/util"
)

const signal = 27      // ASCII 0x1B Escape
const addrMaxLen = 253 // https://en.wikipedia.org/wiki/Domain_Name_System

var readWraper = util.NewIDWapper(addrMaxLen)
var writeWraper = util.NewIDWapper(1 + addrMaxLen)

type TargetConn struct {
	net.Conn
	Addr string
}

// WriteAddr write target address into connection
func WriteAddr(conn net.Conn, addr string) (net.Conn,error) {
	if err:= writeWraper.WriteID(conn, string(signal)+addr);err!=nil{
		conn.Close()
		return nil,err
	}
	return conn,nil
}

// ParseAddr parse target address from the connection
func ParseAddr(conn net.Conn) (_ *TargetConn, err error) {
	teeConn := &util.TeeConn{Conn: conn}
	teeConn.StartOrReset()
	defer func() {
		teeConn.Stop()
		if err != nil {
			conn.Close()
		}
	}()

	buf := make([]byte, 1)
	if n, err := teeConn.Read(buf); err != nil {
		return nil, fmt.Errorf("Read conn fail: %v, readed: %v", err, buf[:n])
	}

	switch buf[0] {
	case signal:
		addr, err := readWraper.ReadID(conn)
		if err != nil {
			return nil, err
		}

		return &TargetConn{conn, addr}, nil

	case 22: // https, SSL handleshake 22(0x16)
		teeConn.StartOrReset()

		host, _, err := extractSNI(io.Reader(teeConn))
		if err != nil {
			return nil, err
		}

		// not able to parse port from tls connection
		return &TargetConn{teeConn, host + ":443"}, nil

	default: // http
		teeConn.StartOrReset()

		resp, err := http.ReadRequest(bufio.NewReader(teeConn))
		if err != nil {
			return nil, err
		}

		if strings.Contains(resp.Host, ":") {
			return &TargetConn{teeConn, resp.Host}, nil
		}
		return &TargetConn{teeConn, resp.Host + ":80"}, nil
	}
}

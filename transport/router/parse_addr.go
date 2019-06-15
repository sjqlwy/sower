package router

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"

	"github.com/pkg/errors"

	"github.com/wweir/sower/util"
)

const signal = 27 // ASCII 0x1B Escape

// WithTarget write target address into connection
func WithTarget(conn net.Conn, targetAddr string) (net.Conn, error) {
	length := len(targetAddr)
	if length == 0 {
		return conn, nil
	} else if length > 255 {
		conn.Close()
		return nil, errors.Errorf("target address(%s) is too long", targetAddr)
	}

	data := make([]byte, 0, 2+length)
	data = append(data, signal, byte(length))
	data = append(data, []byte(targetAddr)...)

	for nn := 0; nn < 2+length; {
		n, err := conn.Write(data)
		if err != nil {
			conn.Close()
			return nil, err
		}
		nn += n
	}
	return conn, nil
}

// ParseAddr parse target address from the connection
func ParseAddr(conn net.Conn) (newConn net.Conn, addr string, err error) {
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
		return nil, "", fmt.Errorf("Read conn fail: %v, readed: %v", err, buf[:n])
	}

	switch buf[0] {
	case signal:
		if _, err := io.ReadFull(conn, buf); err != nil {
			return nil, "", err
		}

		buf = make([]byte, int(buf[0]) /*length*/)
		if _, err := io.ReadFull(conn, buf); err != nil {
			return nil, "", err
		}

		return conn, string(buf), nil

	case 22: // https, SSL handleshake 22(0x16)
		teeConn.StartOrReset()

		host, _, err := extractSNI(io.Reader(teeConn))
		if err != nil {
			return nil, "", err
		}

		// not able to parse port from tls connection
		return teeConn, host + ":443", nil

	default: // http
		teeConn.StartOrReset()

		resp, err := http.ReadRequest(bufio.NewReader(teeConn))
		if err != nil {
			return nil, "", err
		}

		if strings.Contains(resp.Host, ":") {
			return teeConn, resp.Host, nil
		}
		return teeConn, resp.Host + ":80", nil
	}
}

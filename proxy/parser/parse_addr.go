package parser

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

// InitTarget write target address into connection
func InitTarget(conn net.Conn, targetAddr string) error {
	length := len(targetAddr)
	if length > 255 {
		return errors.New("")
	}

	data := make([]byte, 0, 2+length)
	data = append(data, 13, byte(length))
	data = append(data, []byte(targetAddr)...)

	for nn := 0; nn < 2+length; {
		n, err := conn.Write(data)
		if err != nil {
			return err
		}
		nn += n
	}
	return nil
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
	case 13: // [D]irect connection signal 13(0x0d)
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

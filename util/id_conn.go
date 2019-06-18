package util

import (
	"bytes"
	"io"
	"net"

	"github.com/pkg/errors"
)

type wrapID struct {
	maxLength int
}

func NewIDWapper(maxLength int) *wrapID {
	return &wrapID{maxLength: maxLength}
}

func (c *wrapID) WriteID(conn net.Conn, id string) error {
	if idLen := len(id); idLen == 0 {
		return errors.New("id is empty")
	} else if idLen > c.maxLength {
		return errors.Errorf("id(%s) is too long")
	}

	buf := make([]byte, c.maxLength)
	copy(buf, []byte(id))

	var err error
	for n, nn := 0, 0; nn < int(c.maxLength); nn += n {
		if n, err = conn.Write(buf[nn:]); err != nil {
			return errors.Wrap(err, "write id")
		}
	}
	return nil
}

func (c *wrapID) ReadID(conn net.Conn) (id string, err error) {
	buf := make([]byte, c.maxLength)
	if _, err = io.ReadFull(conn, buf); err != nil {
		return "", err
	}

	return string(bytes.TrimRight(buf, string(0))), nil
}

package p2p

import (
	"net"
)

const FROM = "pf"
const TO = "pt"
const ID_MAX_LEN = 64

var idPlaceholder []byte

func init() {
	idPlaceholder = make([]byte, 64)
	for i := range idPlaceholder {
		idPlaceholder[i] = 32 // ASCII SPACE
	}
}

func writeID(conn net.Conn, id string) (err error) {
	buf := make([]byte, ID_MAX_LEN)
	copy(buf[len(id):], idPlaceholder)
	copy(buf, []byte(id))

	for nn, n := 0, 0; nn < ID_MAX_LEN; nn += n {
		if n, err = conn.Write(buf[nn:]); err != nil {
			return err
		}
	}
	return nil
}

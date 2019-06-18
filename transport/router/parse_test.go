package router

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"net"
	"net/http"
	"testing"

	"github.com/wweir/sower/util"
)

func TestParseAddr1(t *testing.T) {
	c1, c2 := net.Pipe()

	go func() {
		req, _ := http.NewRequest("GET", "http://wweir.cc", bytes.NewReader([]byte{1, 2, 3}))
		req.Write(c1)
	}()

	tgtConn, err := ParseAddr(c2)
	if err != nil || tgtConn.Addr != "wweir.cc:80" {
		t.Error(err, tgtConn.Addr)
	}

	req, err := http.ReadRequest(bufio.NewReader(c2))
	if err != nil {
		t.Error(err)
	}

	data, err := ioutil.ReadAll(req.Body)
	if err != nil || len(data) != 3 || data[0] != 1 {
		t.Error(err, data)
	}
}

func TestParseAddr2(t *testing.T) {
	c1, c2 := net.Pipe()

	go func() {
		c1.Write(util.HTTPS.PingMsg("wweir.cc"))
	}()

	tgtConn, err := ParseAddr(c2)

	if err != nil || tgtConn.Addr != "wweir.cc:443" {
		t.Error(err, tgtConn.Addr)
	}
}

func TestParseAddr3(t *testing.T) {
	c1, c2 := net.Pipe()

	go func() {
		WriteAddr(c1, "wweir.cc:8080")
		c1.Write(util.HTTPS.PingMsg("wweir.cc"))
	}()

	tgtConn, err := ParseAddr(c2)

	if err != nil || tgtConn.Addr != "wweir.cc:8080" {
		t.Error(err, tgtConn.Addr)
	}
}

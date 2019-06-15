package proxy

import (
	"net"
	"net/url"
	"strings"

	"github.com/golang/glog"
	"github.com/wweir/sower/conf"
	"github.com/wweir/sower/transport"
)

func StartServer() {
	cfg := conf.GetConf().Transport
	u, err := url.Parse(cfg.SelfURI)
	if err != nil {
		glog.Fatalln(err)
	}

	tran := transport.NewTransport(u.Scheme)
	connCh, err := tran.Listen(u.Host)
	if err != nil {
		glog.Fatalln(err)
	}

	for tgtConn := range connCh {
		secs := strings.Split(tgtConn.TargetAddr, ":")
		if len(secs) != 2 {
			glog.Errorln("invalid client connected, target addr:", tgtConn.TargetAddr)
		}

		switch secs[0] {
		case "from": // p2p
		case "to": // p2p
		}

		go func(tgtConn *transport.TargetConn) {
			rc, err := net.Dial("tcp", tgtConn.TargetAddr)
			if err != nil {
				glog.Warningln("dial remote:", err)
			}
			rc.(*net.TCPConn).SetKeepAlive(true)

			relay(tgtConn.Conn, rc)
		}(tgtConn)
	}
}

package proxy

import (
	"net"
	"net/url"

	"github.com/golang/glog"
	"github.com/wweir/sower/conf"
	"github.com/wweir/sower/dns"
	"github.com/wweir/sower/transport"
)

func StartIntelligentProxy() {
	conf := conf.GetConf()
	if conf.Transport.OutletURI == "" {
		return
	}

	u, err := url.Parse(conf.Transport.OutletURI)
	if err != nil {
		glog.Fatalln(err)
	}
	tran := transport.NewTransport(u.Scheme)

	go dns.StartDNS(conf.Client.DNS, conf.Client.ServeIP,
		conf.Client.Suggest.SuggestLevel, conf.AddSuggestion)

	for _, port := range []string{"80", "443"} {
		go func() {
			addr := net.JoinHostPort(conf.Client.ServeIP, port)
			ln, err := net.Listen("tcp", addr)
			if err != nil {
				glog.Fatalln(err)
			}

			for {
				conn, err := ln.Accept()
				if err != nil {
					glog.Errorln("accept", addr, "fail:", err)
					continue
				}
				conn.(*net.TCPConn).SetKeepAlive(true)

				go func(conn net.Conn) {
					rc, err := tran.Dial(u.Host, "")
					if err != nil {
						glog.Errorln("listen socks5 addr fail:", err)
						conn.Close()
						return
					}

					relay(conn, rc)
				}(conn)
			}
		}()
		glog.Infoln("listening port:", port)
	}
}

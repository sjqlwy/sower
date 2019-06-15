package proxy

import (
	"net"
	"net/url"

	"github.com/golang/glog"
	"github.com/wweir/sower/config"
	"github.com/wweir/sower/transport"
)

func StartDirectProxy() {
	for _, proxy := range config.GetConf().DirectProxys {
		go func(listenAddr, outletURI, tgtAddr string) {
			ln, err := net.Listen("tcp", listenAddr)
			if err != nil {
				glog.Fatalln(err)
			}

			u, err := url.Parse(outletURI)
			if err != nil {
				glog.Fatalln(err)
			}
			tran := transport.NewTransport(u.Scheme)

			for {
				conn, err := ln.Accept()
				if err != nil {
					glog.Errorln("listen socks5 addr fail:", err)
					continue
				}
				conn.(*net.TCPConn).SetKeepAlive(true)

				go func(conn net.Conn) {
					rc, err := tran.Dial(u.Host, tgtAddr)
					if err != nil {
						glog.Errorln("listen socks5 addr fail:", err)
						conn.Close()
						return
					}

					relay(conn, rc)
				}(conn)
			}
		}(proxy.ListenAddr, proxy.OutletURI, proxy.TargetAddr)
	}
}

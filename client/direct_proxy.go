package client

import (
	"net"

	"github.com/golang/glog"
	"github.com/wweir/sower/config"
)

func StartDirectProxy() {
	for _, proxy := range config.GetCfg().DirectProxys {
		ln, err := net.Listen("tcp", proxy.ListenAddr)
		if err != nil {
			glog.Fatalln(err)
		}

		go func(ln net.Listener, tgtAddr, outletURI string) {
			for {
				conn, err := ln.Accept()
				if err != nil {
					glog.Errorln("listen socks5 addr fail:", err)
					return
				}
				conn.(*net.TCPConn).SetKeepAlive(true)

				// if socks5Addr != "" {
				// 	rc, err := socks5.Dial(socks5Addr, tgtAddr)
				// 	if err != nil {
				// 		glog.Errorln("dial socks5 addr fail:", err)
				// 	}
				// 	relay(conn, rc)
				// 	return
				// }

				// TODO: p2p  relay
			}
		}(ln, proxy.TargetAddr, proxy.OutletURI)
	}
}

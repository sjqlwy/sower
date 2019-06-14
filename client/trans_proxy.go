package client

import (
	"net"

	"github.com/golang/glog"
	"github.com/wweir/sower/client/dns"
	"github.com/wweir/sower/config"
	"github.com/wweir/sower/proxy/parser"
	"github.com/wweir/sower/proxy/socks5"
)

func StartClient() {
	cfg := config.GetCfg().Client
	if cfg.ClientIP == "" {
		return
	}

	go dns.StartDNS(cfg.DNSIP, cfg.ClientIP, cfg.Suggest.SuggestLevel, config.GetCfg().AddSuggestion)

	for _, port := range []string{":80", ":443"} {
		go func() {
			ln, err := net.Listen("tcp", cfg.ClientIP+port)
			if err != nil {
				glog.Fatalln(err)
			}

			for {
				conn, err := ln.Accept()
				if err != nil {
					glog.Errorln("accept", cfg.ClientIP+port, "fail:", err)
					continue
				}
				conn.(*net.TCPConn).SetKeepAlive(true)

				if cfg.Socks5Addr != "" {
					if conn, addr, err := parser.ParseAddr(conn); err != nil {
						glog.Errorln(err)
						conn.Close()
						continue

					} else if rc, err := socks5.Dial(cfg.Socks5Addr, addr); err != nil {
						glog.Errorln(err)
						conn.Close()
						continue
					} else {
						relay(conn, rc)
						continue
					}
				}

				// TODO: p2p relay
			}
		}()
		glog.Infoln("listening port:", port)
	}
}

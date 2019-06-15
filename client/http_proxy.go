package client

import (
	"context"
	"crypto/tls"
	"io"
	"net"
	"net/http"
	"time"

	"github.com/golang/glog"
	"github.com/wweir/sower/config"
)

func StartHttpProxy() {
	for _, proxy := range config.GetCfg().HTTPProxys {
		srv := &http.Server{
			Addr: proxy.ListenAddr,
			// Disable HTTP/2.
			TLSNextProto: map[string]func(*http.Server, *tls.Conn, http.Handler){},
			IdleTimeout:  90 * time.Second,
		}
		srv.Handler= http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodConnect {
					httpsProxy(w, r, proxy.OutletURI)
				} else {
					httpProxy(w, r, proxy.OutletURI)
				}
			})

		glog.Infoln("listening http proxy on", proxy.ListenAddr)
		go glog.Fatalln(srv.ListenAndServe())
	}

}

func httpProxy(w http.ResponseWriter, r *http.Request, outletURI string) {
	roundTripper := &http.Transport{
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
		roundTripper.DialContext = func(context.Context, string, string) (net.Conn, error) {
			// FIXME: p2p
			return nil, nil
		}

	resp, err := roundTripper.RoundTrip(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		glog.Errorln("serve https proxy, get remote data:", err)
		return
	}
	defer resp.Body.Close()

	for k, vs := range resp.Header {
		for _, v := range vs {
			w.Header().Add(k, v)
		}
	}
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func httpsProxy(w http.ResponseWriter, r *http.Request,  outletURI string) {
	// local conn
	conn, _, err := w.(http.Hijacker).Hijack()
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	conn.(*net.TCPConn).SetKeepAlive(true)

	if _, err := conn.Write([]byte(r.Proto + " 200 Connection established\r\n\r\n")); err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		conn.Close()
		glog.Errorln("serve https proxy, write data fail:", err)
		return
	}

	// if Socks5Addr != "" {
	// 	rc, err := socks5.Dial(Socks5Addr, r.Host)
	// 	if err != nil {
	// 		http.Error(w, err.Error(), http.StatusServiceUnavailable)
	// 		conn.Close()
	// 		glog.Errorln("serve https proxy, dial remote fail:", err)
	// 		return
	// 	}
	// 	relay(conn, rc)
	// 	return
	// }

	// remote conn
	// host, port, _ := net.SplitHostPort(r.Host)
	// TODO: p2p relay, if port == 443, run as trans, or run as direct
}

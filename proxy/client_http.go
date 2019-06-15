package proxy

import (
	"context"
	"crypto/tls"
	"io"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/golang/glog"
	"github.com/wweir/sower/conf"
	"github.com/wweir/sower/transport"
)

func StartHttpProxy() {
	for _, proxy := range conf.GetConf().HTTPProxys {
		srv := &http.Server{
			Addr: proxy.ListenAddr,
			// Disable HTTP/2.
			TLSNextProto: map[string]func(*http.Server, *tls.Conn, http.Handler){},
			IdleTimeout:  90 * time.Second,
		}

		u, err := url.Parse(proxy.OutletURI)
		if err != nil {
			glog.Fatalln(err)
		}
		tran := transport.NewTransport(u.Scheme)

		srv.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodConnect {
				httpsProxy(w, r, tran, u.Host)
			} else {
				httpProxy(w, r, tran, u.Host)
			}
		})

		glog.Infoln("listening http proxy on", proxy.ListenAddr)
		go glog.Fatalln(srv.ListenAndServe())
	}

}

func httpProxy(w http.ResponseWriter, r *http.Request, tran transport.Transport, addr string) {
	roundTripper := &http.Transport{
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	roundTripper.DialContext = func(context.Context, string, string) (net.Conn, error) {
		return tran.Dial(addr, "")
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

func httpsProxy(w http.ResponseWriter, r *http.Request, tran transport.Transport, addr string) {
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

	var rc net.Conn
	if _, port, err := net.SplitHostPort(r.Host); err == nil && port != "443" {
		rc, err = tran.Dial(addr, r.Host)
	} else {
		rc, err = tran.Dial(addr, "")
	}

	relay(conn, rc)
}

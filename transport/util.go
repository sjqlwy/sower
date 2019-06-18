package transport

import (
	"net"
	"net/url"

	"github.com/pkg/errors"
	"github.com/wweir/sower/transport/router"
)

type Transport interface {
	Dial(addr string) (net.Conn, error)
	Listen(addr string) (<-chan *router.TargetConn, error)
}

// NewCSTransport return a c/s connect mode transport
func NewCSTransport(remoteURI string) (Transport, error) {
	u, err := url.Parse(remoteURI)
	if err != nil {
		return nil, err
	}

	switch u.Scheme {
	case "tcp":
		return NewTCP(u.Host), nil
	case "quic":
		return NewQUIC(u.Host), nil
	case "kcp":
		return NewKCP(u.Host), nil
	default:
		return nil, errors.Errorf("invalid transport type: %s in %s", u.Scheme, u)
	}
}
func NewTransport(remoteURI, localURI, brokerURI string) (Transport, error) {
	if tran, err := NewCSTransport(remoteURI); err == nil {
		return tran, nil
	}

	u, err := url.Parse(remoteURI)
	if err != nil {
		return nil, err
	}
	switch u.Scheme {
	case "socks5", "socks5h":
		return NewSocks5(u.Host), nil
	}

	uLocal, err := url.Parse(remoteURI)
	if err != nil {
		return nil, err
	}
	brokerTran, err := NewCSTransport(remoteURI)
	if err != nil {
		return nil, err
	}

	switch u.Scheme {
	case "tcp_p2p":
		return NewTCPP2P(uLocal.Host, u.Host, brokerTran)
	default:
		return nil, errors.Errorf("invalid transport type: %s in %s", u.Scheme, u)
	}
}

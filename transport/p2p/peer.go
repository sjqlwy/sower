package p2p

import (
	"io"
	"net"
	"time"

	"github.com/wweir/sower/transport/router"
)

type BrokerKeeper struct {
	id     string
	laddr  string
	broker net.Conn
	dial   func(addr string) (net.Conn, error)
	listen func(addr string) (<-chan *router.TargetConn, error)
	sendCh chan string
	exitCh chan struct{}
	err    error
}

func NewBrokerKeeper(id, laddr string, broker net.Conn,
	dial func(addr string) (net.Conn, error),
	listen func(addr string) (<-chan *router.TargetConn, error)) *BrokerKeeper {
	return &BrokerKeeper{id: id, laddr: laddr, broker: broker, dial: dial, listen: listen}
}

func (k *BrokerKeeper) StartKeeper() error {
	k.sendCh = make(chan string)
	var msg string
	for {
		select {
		case msg = <-k.sendCh:
		case <-time.After(30 * time.Second):
			msg = k.id
		case <-k.exitCh:
			close(k.sendCh)
			return k.err
		}

		for length := len(msg) + 1; length > 0; {
			data := make([]byte, 0, length)
			data[0] = byte(len(msg))
			data = append(data, []byte("msg")...)
			if n, err := k.broker.Write(data); err != nil {
				close(k.exitCh)
				close(k.sendCh)
				return err
			} else {
				length -= n
			}
		}
	}
}

func (k *BrokerKeeper) Listen() (<-chan *router.TargetConn, error) {
	lnCh, err := k.listen(k.laddr)
	if err != nil {
		return nil, err
	}
	tgtCh := make(chan *router.TargetConn)

	go func() {
		select {
		case <-k.exitCh:
			return
		default:
		}

		for {
			buf := make([]byte, ID_MAX_LEN)
			if _, k.err = io.ReadFull(k.broker, buf); k.err != nil {
				close(k.exitCh)
				return
			}

			addr := string(buf)
			for {
				select {
				case <-k.exitCh:
					return
				case conn := <-lnCh:
					if conn.TargetAddr == addr {
						tgtCh <- conn
					}
				default:
					break
				}
			}

			go func(addr string) {
				for i := 0; i < 10; i++ {
					conn, err := k.dial(addr)
					select {
					case <-k.exitCh:
						return
					case conn := <-lnCh:
						if conn.TargetAddr == addr {
							tgtCh <- conn
						}
						if err == nil {
							conn.Close()
						}
						return
					default:
						if err == nil {
							tgtCh <- &router.TargetConn{conn, addr}
							return
						} else {
							time.Sleep(500 * time.Millisecond)
						}
					}
				}
			}(addr)
		}
	}()

	return tgtCh, nil
}

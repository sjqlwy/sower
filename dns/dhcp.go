package dns

import (
	"math/rand"
	"net"
	"runtime"
	"time"

	"github.com/krolaw/dhcp4"
	"github.com/libp2p/go-reuseport"
	"github.com/pkg/errors"
	"github.com/wweir/sower/util"
)

var xid = make([]byte, 4)
var broadcastAddr, _ = net.ResolveUDPAddr("udp", "255.255.255.255:67")

// GetDefaultDNSServer return default dns server with dhcpv4 protocol
func GetDefaultDNSServer() (ip, dns net.IP, err error) {
	iface, err := util.PickInterface()
	if err != nil {
		return nil, nil, errors.Wrap(err, "pick interface")
	}

	rand.Read(xid)
	pack := dhcp4.RequestPacket(dhcp4.Discover, iface.HardwareAddr,
		net.IPv4(0, 0, 0, 0), xid, true, []dhcp4.Option{
			{Code: dhcp4.OptionRequestedIPAddress, Value: []byte(iface.IP.To4())},
			{Code: dhcp4.End},
		})

	var conn net.PacketConn
	if runtime.GOOS == "windows" {
		if conn, err = reuseport.ListenPacket("udp4", iface.IP.String()+":68"); err != nil {
			return nil, nil, errors.Wrap(err, "listen dhcp")
		}
	} else {
		if conn, err = reuseport.ListenPacket("udp4", "0.0.0.0:68"); err != nil {
			return nil, nil, errors.Wrap(err, "listen dhcp")
		}
	}
	defer conn.Close()
	conn.SetDeadline(time.Now().Add(time.Second))

	if _, err := conn.WriteTo([]byte(pack), broadcastAddr); err != nil {
		return nil, nil, errors.Wrap(err, "write broadcast")
	}

	buf := make([]byte, 1500 /*MTU*/)
	n, _, err := conn.ReadFrom(buf)
	if err != nil {
		return nil, nil, errors.Wrap(err, "read dhcp offer")
	}

	options := dhcp4.Packet(buf[:n]).ParseOptions()
	ipBytes := options[dhcp4.OptionServerIdentifier]
	dnsBytes := options[dhcp4.OptionDomainNameServer]
	return net.IP(ipBytes), net.IP(dnsBytes), nil
}

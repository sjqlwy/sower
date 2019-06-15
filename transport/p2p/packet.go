package p2p

import (
	"github.com/google/uuid"
)

type header struct {
	// (proto version)<<4 + (net version)
	//	proto version:
	//		0x1. => the current verion
	//	net version:
	//		0x.1 => IPv4
	//		0x.2 => IPv6
	version byte

	// command action, the highest 4 bit means command kind
	// A/B are peers, A connect to B, S is broker server
	//	0x1. => peer to server
	//		0x11 => P connect to S
	//		0x12 => A ask S wanna connect to B
	//		0x13 => B tell S, connect to A result
	//	0x2. => server to peer
	//		0x21 => S heartbeat to P
	//		0x22 => S tell A, B connect to A result
	//		0x23 => S ask B, try connect to A
	//	0x3. => peer to peer
	//		0x31 => B try connect to A, package mostly likely to be dropped
	//		0x32 => A connect to B
	command byte
}

const (
	v1IPv4 = 0x11 + iota
	v1IPv6
)

const (
	p2s = 0x11 + iota
	a2s
	b2s
)
const (
	s2p = 0x21 + iota
	s2a
	s2b
)
const (
	b2a = 0x31 + iota
	a2b
)

type p2sT struct {
	uuid uuid.UUID
}
type a2sT struct {
	uuid uuid.UUID
}

type s2bTv4 struct {
	ip   [4]byte
	port uint16
}
type s2bTv6 struct {
	ip   [16]byte
	port uint16
}
type s2aT struct {
}

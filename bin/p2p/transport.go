package p2p

import "net"

// Peer represents a connected remote peer
type Peer interface {
	Conn() net.Conn
	Close() error
}

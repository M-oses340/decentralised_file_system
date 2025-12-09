package p2p

import "net"

type TCPPeer struct {
	conn     net.Conn
	outbound bool
}

// Constructor used by TCPTransport
func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer{
		conn:     conn,
		outbound: outbound,
	}
}

func (p *TCPPeer) Conn() net.Conn {
	return p.conn
}

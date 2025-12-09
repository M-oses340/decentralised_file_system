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

// Conn returns the underlying net.Conn
func (p *TCPPeer) Conn() net.Conn {
	return p.conn
}

// Close closes the connection
func (p *TCPPeer) Close() error {
	return p.conn.Close()
}

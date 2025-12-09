package p2p

import (
	"fmt"
	"net"
	"sync"
)

type TCPPeer struct {
	conn     net.Conn
	outbound bool
}

func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer{conn: conn, outbound: outbound}
}

type TCPTransport struct {
	listenAddr string
	listener   net.Listener

	mu    sync.RWMutex
	peers map[net.Addr]Peer
}

func NewTCPTransport(listenAddr string) *TCPTransport {
	return &TCPTransport{
		listenAddr: listenAddr,
		peers:      make(map[net.Addr]Peer), // important!
	}
}

func (t *TCPTransport) ListenAndAccept() error {
	var err error

	// ❌ WRONG: t.listener err = ...
	// ✔️ FIX:
	t.listener, err = net.Listen("tcp", t.listenAddr)
	if err != nil {
		return err
	}

	go t.startAcceptLoop()
	return nil
}

func (t *TCPTransport) startAcceptLoop() {
	for {
		conn, err := t.listener.Accept()
		if err != nil {
			fmt.Printf("TCP accept error: %s\n", err)
			continue
		}

		// ❌ WRONG (placed outside loop)
		// ✔️ FIX: Handle each connection inside loop
		go t.handleConn(conn)
	}
}

func (t *TCPTransport) handleConn(conn net.Conn) {
	peer := NewTCPPeer(conn, true)
	fmt.Println("New incoming connection from", peer)

	// TODO: wrap net.Conn into your Peer object
	// Example:
	// peer := NewPeer(conn)
	// t.addPeer(peer)

}

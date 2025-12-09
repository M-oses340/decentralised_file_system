package p2p

import (
	"fmt"
	"net"
	"sync"
)

type TCPTransport struct {
	listenAddr    string
	listener      net.Listener
	HandshakeFunc HandshakeFunc
	decoder       Decoder

	mu    sync.RWMutex
	peers map[net.Addr]Peer
}

func NewTCPTransport(listenAddr string) *TCPTransport {
	return &TCPTransport{
		listenAddr:    listenAddr,
		peers:         make(map[net.Addr]Peer),
		HandshakeFunc: NOPHandshakeFunc(),
		decoder:       NewDefaultDecoder(),
	}
}

func (t *TCPTransport) ListenAndAccept() error {
	var err error

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

		go t.handleConn(conn)
	}
}

type Temp struct{}

func (t *TCPTransport) handleConn(conn net.Conn) {
	peer := NewTCPPeer(conn, false)

	if err := t.HandshakeFunc(peer); err != nil {
		fmt.Printf("Handshake failed: %s\n", err)
		conn.Close()
		return
	}

	t.mu.Lock()
	t.peers[conn.RemoteAddr()] = peer
	t.mu.Unlock()

	fmt.Printf("New peer connected: %s\n", conn.RemoteAddr())

	msg := &Temp{}

	for {
		if err := t.decoder.Decode(conn, msg); err != nil {
			fmt.Printf("TCP decode error: %s\n", err)
			return // stop loop when decode fails
		}

		fmt.Printf("Received message: %+v\n", msg)
	}
}

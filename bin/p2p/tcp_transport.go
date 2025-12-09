package p2p

import (
	"fmt"
	"net"
	"sync"
)

// TCPTransportOpts defines options for TCPTransport
type TCPTransportOpts struct {
	ListenAddr    string
	HandshakeFunc HandshakeFunc
	Decoder       Decoder
	OnPeer        func(Peer) error
}

// TCPTransport manages multiple TCP peers
type TCPTransport struct {
	listenAddr    string
	listener      net.Listener
	HandshakeFunc HandshakeFunc
	decoder       Decoder
	OnPeer        func(Peer) error

	mu    sync.RWMutex
	peers map[net.Addr]Peer
}

// Temp is a simple message type
type Temp struct {
	Value string
}

// NewTCPTransport creates a TCPTransport
func NewTCPTransport(opts TCPTransportOpts) *TCPTransport {
	if opts.HandshakeFunc == nil {
		opts.HandshakeFunc = NOPHandshakeFunc()
	}
	if opts.Decoder == nil {
		opts.Decoder = NewDefaultDecoder()
	}

	return &TCPTransport{
		listenAddr:    opts.ListenAddr,
		HandshakeFunc: opts.HandshakeFunc,
		decoder:       opts.Decoder,
		OnPeer:        opts.OnPeer,
		peers:         make(map[net.Addr]Peer),
	}
}

// ListenAndAccept starts listening and accepting connections
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

func (t *TCPTransport) handleConn(conn net.Conn) {
	peer := NewTCPPeer(conn, false)

	// handshake
	if err := t.HandshakeFunc(peer); err != nil {
		fmt.Printf("Handshake error: %s\n", err)
		conn.Close()
		return
	}

	// register peer
	t.mu.Lock()
	t.peers[conn.RemoteAddr()] = peer
	t.mu.Unlock()

	fmt.Printf("New peer connected: %s\n", conn.RemoteAddr())

	// optional callback
	if t.OnPeer != nil {
		if err := t.OnPeer(peer); err != nil {
			fmt.Printf("OnPeer error: %s\n", err)
		}
	}

	// read messages
	for {
		msg := &Temp{}
		if err := t.decoder.Decode(conn, msg); err != nil {
			fmt.Printf("TCP decode error: %s\n", err)
			return
		}
		fmt.Printf("Received message: %+v\n", msg)
	}
}

package p2p

import (
	"fmt"
	"net"
	"sync"
)

type TCPTransportOpts struct {
	ListenAddr    string
	HandshakeFunc HandshakeFunc
	Decoder       Decoder
	OnPeer        func(Peer) error
}

type TCPTransport struct {
	listenAddr    string
	listener      net.Listener
	HandshakeFunc HandshakeFunc
	decoder       Decoder
	OnPeer        func(Peer) error

	mu    sync.RWMutex
	peers map[net.Addr]Peer
}

func NewTCPTransport(opts TCPTransportOpts) *TCPTransport {

	// Apply defaults if options missing
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

	// Run handshake
	if err := t.HandshakeFunc(peer); err != nil {
		fmt.Printf("Handshake failed: %s\n", err)
		conn.Close()
		return
	}

	// Add peer
	t.mu.Lock()
	t.peers[conn.RemoteAddr()] = peer
	t.mu.Unlock()

	fmt.Printf("New peer connected: %s\n", conn.RemoteAddr())

	// Fire callback if provided
	if t.OnPeer != nil {
		if err := t.OnPeer(peer); err != nil {
			fmt.Printf("OnPeer error: %s\n", err)
		}
	}

	// Read messages forever
	msg := &Temp{}
	for {
		if err := t.decoder.Decode(conn, msg); err != nil {
			fmt.Printf("TCP decode error: %s\n", err)
			return
		}
		fmt.Printf("Received message: %+v\n", msg)
	}
}

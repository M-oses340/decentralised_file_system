package p2p

// HandshakeFunc runs immediately after peer connects
type HandshakeFunc func(Peer) error

// NOPHandshakeFunc returns a handshake function that does nothing
func NOPHandshakeFunc() HandshakeFunc {
	return func(p Peer) error {
		return nil
	}
}

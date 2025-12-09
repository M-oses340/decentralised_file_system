package p2p

type HandshakeFunc func(Peer) error

func NOPHandshakeFunc() HandshakeFunc {
	return func(_ Peer) error {
		return nil
	}
}

package p2p

import (
	"net"
)

type Decoder interface {
	Decode(conn net.Conn, v any) error
}

type DefaultDecoder struct{}

func NewDefaultDecoder() Decoder {
	return &DefaultDecoder{}
}

func (d *DefaultDecoder) Decode(conn net.Conn, v any) error {
	buf := make([]byte, 256)
	_, err := conn.Read(buf)
	return err
}

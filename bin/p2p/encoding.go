package p2p

import (
	"bufio"
	"fmt"
	"io"
)

// Decoder reads a message from a connection into v
type Decoder interface {
	Decode(r io.Reader, v any) error
}

// DefaultDecoder reads line-based messages
type DefaultDecoder struct{}

func NewDefaultDecoder() Decoder {
	return &DefaultDecoder{}
}

func (d *DefaultDecoder) Decode(r io.Reader, v any) error {
	reader := bufio.NewReader(r)
	line, err := reader.ReadString('\n')
	if err != nil {
		return err
	}

	switch m := v.(type) {
	case *Temp:
		m.Value = line
	default:
		return fmt.Errorf("unknown message type: %T", v)
	}

	return nil
}

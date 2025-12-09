package p2p

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTransport(t *testing.T) {
	listenAddr := "127.0.0.1:4000" // fixed: explicit IP is safer
	tr := NewTCPTransport(listenAddr)

	assert.Equal(t, listenAddr, tr.listenAddr)

	// Start the listener
	err := tr.ListenAndAccept()
	assert.Nil(t, err)

	// Give the goroutine a moment
	time.Sleep(50 * time.Millisecond)

	// Verify that the listener is created
	assert.NotNil(t, tr.listener)

	// Cleanup: close the listener after the test
	err = tr.listener.Close()
	assert.Nil(t, err)

	select {}
}

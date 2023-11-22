package tcpserver

import (
	"net"
	"testing"
	"time"
)

func TestServerInitialization(t *testing.T) {
	// Set up the server
	addr := "localhost:8080"
	server := New(addr)

	// Start server in a goroutine so that it doesn't block
	go func() {
		err := server.Start()
		if err != nil {
			t.Errorf("failed to start server: %s", err)
		}
	}()

	// We wait for the server to start
	time.Sleep(2 * time.Second)

	// Try to stablish a connection with the server
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		t.Errorf("failed to connect to server: %s", err)
	}
	conn.Close()
}

package tcpserver

import (
	"fmt"
	"net"
)

type Client struct {
	conn net.Conn
}

func NewClient(connection net.Conn) *Client {
	return &Client{conn: connection}
}

func (c *Client) HandleConnection() {
	defer c.conn.Close() // Ensure the connection is closed

	fmt.Printf("Accepted new client: %s\n", c.conn.RemoteAddr())

	// Logic to handle client communication goes here
}

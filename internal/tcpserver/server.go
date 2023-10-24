package tcpserver

import (
	"fmt"
	"net"
)

type Server struct {
	address string
}

func New(address string) *Server {
	return &Server{address: address}
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.address)
	if err != nil {
		return fmt.Errorf("error setting up server listener: %w", err)
	}
	defer ln.Close()

	fmt.Printf("TCP Sever started on %s\n", s.address)

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Printf("Error accepting connection: %v\n", err)
			continue
		}

		client := NewClient(conn)
		go client.HandleConnection()

	}
}

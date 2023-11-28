package tcpserver

import (
	"fmt"
	"net"
	"sync"
)

// Message represents a simple structure to hold date from a network connection.
type Message struct {
	From    string // Sender's information
	Payload []byte // The actual message data
}

type Server struct {
	address string
	ln      net.Listener
	quitch  chan struct{}
	msgch   chan Message

	// Connection management
	maxConnections    int
	activeConnections int
	connLock          sync.Mutex // Protects activeConnections
	connCond          *sync.Cond // Used to wait for free connection slots
}

func New(address string, maxConnections int) *Server {
	server := &Server{
		address:        address,
		quitch:         make(chan struct{}),
		msgch:          make(chan Message, 10),
		maxConnections: maxConnections,
	}
	server.connCond = sync.NewCond(&server.connLock)
	return &Server{address: address}
}

func (s *Server) Start() error {
	var err error
	s.ln, err = net.Listen("tcp", s.address)
	if err != nil {
		return fmt.Errorf("error setting up server listener: %w", err)
	}

	fmt.Printf("TCP Sever started on %s\n", s.address)

	go s.acceptLoop()

	<-s.quitch // Wait for a quit signal before shutting down
	return nil
}

// acceptLoop handles incoming connections.
func (s *Server) acceptLoop() {
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			fmt.Printf("error accepting connection: %v\n", err)
			continue
		}

		// check if we've reached the maximum number of connections
		s.connLock.Lock()
		if s.activeConnections >= s.maxConnections {
			s.connLock.Unlock()
			conn.Close() // Refuse new connection
			continue
		}
		// Register the new connection
		s.activeConnections++
		s.connLock.Unlock()

		// Handle the connection in a separate goroutine
		go s.handleConnection(conn)
	}
}

// handleConnection processes a single client connection.
func (s *Server) handleConnection(conn net.Conn) {
	defer func() {
		conn.Close()
		s.connLock.Lock()
		s.activeConnections--
		s.connCond.Broadcast() // Notify that a connection slot is free
		s.connLock.Unlock()
	}()
	// Here, you'd handle your connection, such as reading message, etc.
}

package tcpserver

import (
	"fmt"
	"io"
	"log"
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
	connLock          sync.Mutex     // Protects activeConnections
	connCond          *sync.Cond     // Used to wait for free connection slots
	shutdown          chan struct{}  // Channel to signal server shutdown
	wg                sync.WaitGroup // WaitGroup to wait for goroutines to finish
}

func New(address string, maxConnections int) *Server {
	server := &Server{
		address:        address,
		quitch:         make(chan struct{}),
		msgch:          make(chan Message, 10),
		maxConnections: maxConnections,
		shutdown:       make(chan struct{}),
	}
	server.connCond = sync.NewCond(&server.connLock)
	return server
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
			select {
			case <-s.shutdown:
				return // Server is shutting down
			default:
				fmt.Printf("error accepting connection: %v\n", err)
			}
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
		log.Print("New connection accepted")

		// Handle the connection in a separate goroutine
		go s.handleConnection(conn)
	}
}

// handleConnection processes a single client connection.
func (s *Server) handleConnection(conn net.Conn) {
	s.wg.Add(1)       // Increment the waitGrooup counter
	defer s.wg.Done() // Decrement the counter when the goroutine completes
	defer func() {
		conn.Close()
		s.connLock.Lock()
		s.activeConnections--
		s.connCond.Broadcast() // Notify that a connection slot is free
		s.connLock.Unlock()
	}()

	buf := make([]byte, 1024) // Buffer for reading data

	for {
		// Read from the connection
		n, err := conn.Read(buf)
		if err != nil {
			if err != io.EOF {
				fmt.Printf("Read error: %v\n", err)
			}
			break
		}
		data := buf[:n]

		// Echo de data back to the client
		_, err = conn.Write(data)
		if err != nil {
			fmt.Printf("Write error: %v\n", err)
			break // Failed to write to the connection
		}
	}
}

func (s *Server) Shutdown() {
	close(s.shutdown) // Signal all goroutines to stop
	s.ln.Close()      // Close the listener

	s.wg.Wait() // wait for all goroutines to finish

	close(s.msgch)  // Close the message channel
	close(s.quitch) //
}

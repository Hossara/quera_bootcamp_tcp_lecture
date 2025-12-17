package server

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
	"tcp_server/protocol"
	"time"
)

// Server represents a TCP server that handles multiple clients
type Server struct {
	address  string               // Address to listen on (e.g., ":8080")
	listener net.Listener         // TCP listener
	clients  map[net.Conn]*Client // Connected clients
	messages []StoredMessage      // Message history
	mu       sync.RWMutex         // Mutex for thread-safe client map access
	quit     chan struct{}        // Channel to signal server shutdown
}

// StoredMessage represents a stored chat message
type StoredMessage struct {
	From      string    // Username of sender
	Content   string    // Message content
	Timestamp time.Time // When the message was sent
}

// Client represents a connected client with metadata
type Client struct {
	conn     net.Conn   // TCP connection
	username string     // Client's username (set via REGISTER command)
	mu       sync.Mutex // Mutex for thread-safe client access
}

// NewServer creates a new TCP server
func NewServer(address string) *Server {
	return &Server{
		address:  address,
		clients:  make(map[net.Conn]*Client),
		messages: make([]StoredMessage, 0, 100), // Preallocate for 100 messages
		quit:     make(chan struct{}),
	}
}

// Start begins listening for connections and accepts clients
func (s *Server) Start() error {
	// Create TCP listener
	listener, err := net.Listen("tcp", s.address)
	if err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}
	s.listener = listener

	log.Printf("🚀 TCP Server started on %s", s.address)
	log.Println("📖 Educational TCP/IP Server - Ready to accept connections")
	log.Println("-----------------------------------------------------------")

	// Accept connections in a loop
	go s.acceptConnections()

	// Wait for shutdown signal
	<-s.quit
	return nil
}

// acceptConnections continuously accepts new client connections
func (s *Server) acceptConnections() {
	for {
		// Accept() blocks until a new connection arrives
		conn, err := s.listener.Accept()
		if err != nil {
			select {
			case <-s.quit:
				// Server is shutting down
				return
			default:
				log.Printf("❌ Error accepting connection: %v", err)
				continue
			}
		}

		log.Printf("✅ New connection from %s", conn.RemoteAddr())

		// Create client object
		client := &Client{
			conn:     conn,
			username: "anonymous",
		}

		// Add client to map
		s.mu.Lock()
		s.clients[conn] = client
		s.mu.Unlock()

		// Handle client in a separate goroutine (concurrent handling)
		// This is KEY: each client gets their own goroutine
		go s.handleClient(client)
	}
}

// handleClient processes messages from a single client
func (s *Server) handleClient(client *Client) {
	defer func() {
		// Cleanup when client disconnects
		log.Printf("👋 Client %s (%s) disconnected", client.conn.RemoteAddr(), client.username)
		client.conn.Close()

		s.mu.Lock()
		delete(s.clients, client.conn)
		s.mu.Unlock()
	}()

	// Create buffered reader for efficient reading
	reader := bufio.NewReader(client.conn)

	// Send welcome message
	welcome := protocol.NewResponse(true, "Welcome to TCP/IP Educational Server!", "")
	s.sendResponse(client, welcome)

	// Read messages in a loop
	for {
		// Set read deadline to detect dead connections
		client.conn.SetReadDeadline(time.Now().Add(5 * time.Minute))

		// Read until newline delimiter
		line, err := reader.ReadString('\n')
		if err != nil {
			// Connection closed or error
			if err.Error() != "EOF" {
				log.Printf("⚠️  Error reading from %s: %v", client.conn.RemoteAddr(), err)
			}
			return
		}

		// Parse message
		var msg protocol.Message
		if v := msg.FromJSON([]byte(line)); err != nil {
			log.Printf("❌ Invalid message from %s: %v", client.conn.RemoteAddr(), v)
			response := protocol.NewResponse(false, "Invalid message format", "")
			s.sendResponse(client, response)
			continue
		}

		// Validate message
		if err := msg.Validate(); err != nil {
			log.Printf("❌ Invalid message from %s: %v", client.conn.RemoteAddr(), err)
			response := protocol.NewResponse(false, fmt.Sprintf("Validation error: %v", err), "")
			s.sendResponse(client, response)
			continue
		}

		log.Printf("📨 Received from %s (%s): %s", client.conn.RemoteAddr(), client.username, msg.String())

		// Process the command
		response := s.processCommand(client, &msg)
		s.sendResponse(client, response)

		// Handle QUIT command
		if msg.Command == protocol.CmdQuit {
			return
		}
	}
}

// processCommand handles different command types
func (s *Server) processCommand(client *Client, msg *protocol.Message) *protocol.Response {
	switch msg.Command {
	case protocol.CmdEcho:
		// Simple echo - send back the data
		return protocol.NewResponse(true, "Echo response", msg.Data)

	case protocol.CmdRegister:
		// Register username
		client.mu.Lock()
		client.username = msg.Data
		client.mu.Unlock()
		log.Printf("✏️  Client %s registered as '%s'", client.conn.RemoteAddr(), msg.Data)
		return protocol.NewResponse(true, fmt.Sprintf("Registration successful. Welcome, %s!", msg.Data), "")

	case protocol.CmdMessage:
		// Broadcast message to all clients
		client.mu.Lock()
		msg.From = client.username
		client.mu.Unlock()

		// Store message in history
		s.storeMessage(msg.From, msg.Data)
		return protocol.NewResponse(true, "Message broadcasted", "")

	case protocol.CmdListUsers:
		// List all connected users
		users := s.getConnectedUsers()
		return protocol.NewResponse(true, "Online users", strings.Join(users, ", "))

	case protocol.CmdListMessages:
		// List recent messages
		messages := s.getRecentMessages(20) // Get last 20 messages
		return protocol.NewResponse(true, "Recent messages", messages)

	case protocol.CmdTime:
		// Return server time
		serverTime := time.Now().Format(time.RFC3339)
		return protocol.NewResponse(true, "Server time", serverTime)

	case protocol.CmdQuit:
		// Client wants to disconnect
		return protocol.NewResponse(true, "Goodbye!", "")

	default:
		return protocol.NewResponse(false, "Unknown command", "")
	}
}

// sendResponse sends a response to a client
func (s *Server) sendResponse(client *Client, response *protocol.Response) {
	data, err := response.ToJSON()
	if err != nil {
		log.Printf("❌ Error marshaling response: %v", err)
		return
	}

	_, err = client.conn.Write(data)
	if err != nil {
		log.Printf("❌ Error sending response to %s: %v", client.conn.RemoteAddr(), err)
	}
}

// getConnectedUsers returns a list of connected usernames
func (s *Server) getConnectedUsers() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	users := make([]string, 0, len(s.clients))
	for _, client := range s.clients {
		client.mu.Lock()
		users = append(users, client.username)
		client.mu.Unlock()
	}
	return users
}

// storeMessage stores a message in the server's history
func (s *Server) storeMessage(from, content string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	msg := StoredMessage{
		From:      from,
		Content:   content,
		Timestamp: time.Now(),
	}

	s.messages = append(s.messages, msg)

	// Keep only last 100 messages to prevent unlimited growth
	if len(s.messages) > 100 {
		s.messages = s.messages[len(s.messages)-100:]
	}

	log.Printf("💾 Stored message from %s (total: %d)", from, len(s.messages))
}

// getRecentMessages returns the last N messages formatted as a string
func (s *Server) getRecentMessages(count int) string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if len(s.messages) == 0 {
		return "No messages yet"
	}

	// Get the last N messages
	start := 0
	if len(s.messages) > count {
		start = len(s.messages) - count
	}

	recentMessages := s.messages[start:]

	// Format messages as a string
	var result strings.Builder
	for i, msg := range recentMessages {
		if i > 0 {
			result.WriteString("\n")
		}
		result.WriteString(fmt.Sprintf("[%s] %s: %s",
			msg.Timestamp.Format("15:04:05"),
			msg.From,
			msg.Content))
	}

	return result.String()
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown() {
	log.Println("🛑 Shutting down server...")
	close(s.quit)

	if s.listener != nil {
		s.listener.Close()
	}

	// Close all client connections
	s.mu.Lock()
	for conn := range s.clients {
		conn.Close()
	}
	s.mu.Unlock()

	log.Println("✅ Server shutdown complete")
}

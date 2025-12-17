package client

import (
	"bufio"
	"fmt"
	"net"
	"sync"
	"tcp_server/protocol"
	"time"
)

// Client represents a TCP client that connects to a server
type Client struct {
	address  string        // Server address to connect to (e.g., "localhost:8080")
	conn     net.Conn      // TCP connection
	reader   *bufio.Reader // Buffered reader for efficient reading
	writer   *bufio.Writer // Buffered writer for efficient writing
	mu       sync.Mutex    // Mutex for thread-safe operations
	username string        // Client's username
}

// NewClient creates a new TCP client
func NewClient(address string) *Client {
	return &Client{
		address: address,
	}
}

// Connect establishes a connection to the server
func (c *Client) Connect() error {
	// Dial creates a TCP connection to the server
	conn, err := net.Dial("tcp", c.address)
	if err != nil {
		return fmt.Errorf("failed to connect to server: %w", err)
	}

	c.conn = conn
	c.reader = bufio.NewReader(conn)
	c.writer = bufio.NewWriter(conn)

	// Read welcome message
	response, err := c.readResponse()
	if err != nil {
		return fmt.Errorf("failed to read welcome message: %w", err)
	}

	fmt.Printf("🎉 %s\n", response.Message)
	return nil
}

// SendMessage sends a message to the server and waits for response
func (c *Client) SendMessage(command, data string) (*protocol.Response, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Create message
	msg := protocol.NewMessage(c.username, command, data)

	// Validate message before sending
	if err := msg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid message: %w", err)
	}

	// Convert to JSON
	jsonData, err := msg.ToJSON()
	if err != nil {
		return nil, fmt.Errorf("failed to marshal message: %w", err)
	}

	// Set write deadline
	c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))

	// Write to connection
	_, err = c.writer.Write(jsonData)
	if err != nil {
		return nil, fmt.Errorf("failed to send message: %w", err)
	}

	// Flush buffer to ensure data is sent
	err = c.writer.Flush()
	if err != nil {
		return nil, fmt.Errorf("failed to flush buffer: %w", err)
	}

	// Wait for response
	response, err := c.readResponse()
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	return response, nil
}

// readResponse reads and parses a response from the server
func (c *Client) readResponse() (*protocol.Response, error) {
	// Set read deadline
	c.conn.SetReadDeadline(time.Now().Add(30 * time.Second))

	// Read until newline
	line, err := c.reader.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Parse response
	var response protocol.Response
	if err := response.FromJSON([]byte(line)); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &response, nil
}

// Register registers a username with the server
func (c *Client) Register(username string) error {
	response, err := c.SendMessage(protocol.CmdRegister, username)
	if err != nil {
		return err
	}

	if !response.Success {
		return fmt.Errorf("registration failed: %s", response.Message)
	}

	c.username = username
	fmt.Printf("✅ %s\n", response.Message)
	return nil
}

// Echo sends an echo request
func (c *Client) Echo(data string) error {
	response, err := c.SendMessage(protocol.CmdEcho, data)
	if err != nil {
		return err
	}

	if response.Success {
		fmt.Printf("📢 Echo: %s\n", response.Data)
	} else {
		fmt.Printf("❌ Error: %s\n", response.Message)
	}
	return nil
}

// SendChatMessage sends a chat message
func (c *Client) SendChatMessage(message string) error {
	response, err := c.SendMessage(protocol.CmdMessage, message)
	if err != nil {
		return err
	}

	if response.Success {
		fmt.Printf("✅ %s\n", response.Message)
	} else {
		fmt.Printf("❌ Error: %s\n", response.Message)
	}
	return nil
}

// ListUsers requests the list of online users
func (c *Client) ListUsers() error {
	response, err := c.SendMessage(protocol.CmdListUsers, "")
	if err != nil {
		return err
	}

	if response.Success {
		fmt.Printf("👥 %s: %s\n", response.Message, response.Data)
	} else {
		fmt.Printf("❌ Error: %s\n", response.Message)
	}
	return nil
}

// ListMessages requests the list of recent messages
func (c *Client) ListMessages() error {
	response, err := c.SendMessage(protocol.CmdListMessages, "")
	if err != nil {
		return err
	}

	if response.Success {
		fmt.Printf("💬 %s:\n%s\n", response.Message, response.Data)
	} else {
		fmt.Printf("❌ Error: %s\n", response.Message)
	}
	return nil
}

// GetServerTime requests the server time
func (c *Client) GetServerTime() error {
	response, err := c.SendMessage(protocol.CmdTime, "")
	if err != nil {
		return err
	}

	if response.Success {
		fmt.Printf("🕐 %s: %s\n", response.Message, response.Data)
	} else {
		fmt.Printf("❌ Error: %s\n", response.Message)
	}
	return nil
}

// Quit sends a quit message and closes the connection
func (c *Client) Quit() error {
	response, err := c.SendMessage(protocol.CmdQuit, "")
	if err != nil {
		return err
	}

	if response.Success {
		fmt.Printf("👋 %s\n", response.Message)
	}

	return c.Close()
}

// StartListening starts listening for broadcast messages from server
// This should be run in a separate goroutine
func (c *Client) StartListening(done chan struct{}) {
	for {
		select {
		case <-done:
			return
		default:
			// Try to read broadcast messages (non-blocking)
			c.conn.SetReadDeadline(time.Now().Add(100 * time.Millisecond))

			line, err := c.reader.ReadString('\n')
			if err != nil {
				// Timeout or error - continue
				continue
			}

			// Parse response
			var response protocol.Response
			if err := response.FromJSON([]byte(line)); err != nil {
				continue
			}

			// Display broadcast message
			if response.Success && len(response.Message) > 0 {
				fmt.Printf("\n%s\n", response.Message)
				fmt.Print("> ")
			}
		}
	}
}

// Close closes the connection
func (c *Client) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// GetUsername returns the client's username
func (c *Client) GetUsername() string {
	return c.username
}

// IsConnected returns true if the client is connected
func (c *Client) IsConnected() bool {
	return c.conn != nil
}

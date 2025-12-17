package protocol

import (
	"encoding/json"
	"fmt"
)

// Message represents a structured message exchanged between client and server
// This demonstrates how to create a protocol for TCP communication
type Message struct {
	Command string `json:"command"` // The action to perform (ECHO, REGISTER, MESSAGE, etc.)
	Data    string `json:"data"`    // The payload/content of the message
	From    string `json:"from"`    // Sender identifier (populated by server)
}

// Response represents the server's response to a client request
type Response struct {
	Success bool   `json:"success"` // Whether the operation succeeded
	Message string `json:"message"` // Response message or error description
	Data    string `json:"data"`    // Optional response data
}

// Command constants - these define the protocol's vocabulary
const (
	CmdEcho         = "ECHO"          // Echo back the data (for testing)
	CmdRegister     = "REGISTER"      // Register a username
	CmdMessage      = "MESSAGE"       // Send a chat message
	CmdListUsers    = "LIST_USERS"    // Get list of online users
	CmdListMessages = "LIST_MESSAGES" // Get list of recent messages
	CmdTime         = "TIME"          // Get server time
	CmdQuit         = "QUIT"          // Disconnect from server
)

// NewMessage creates a new message with the given command and data
func NewMessage(from, command, data string) *Message {
	return &Message{
		From:    from,
		Command: command,
		Data:    data,
	}
}

// NewResponse creates a new response
func NewResponse(success bool, message, data string) *Response {
	return &Response{
		Success: success,
		Message: message,
		Data:    data,
	}
}

// ToJSON converts a message to JSON bytes
func (m *Message) ToJSON() ([]byte, error) {
	data, err := json.Marshal(m)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal message: %w", err)
	}
	// Add newline delimiter for easy reading
	return append(data, '\n'), nil
}

// FromJSON parses a message from JSON bytes
func (m *Message) FromJSON(data []byte) error {
	if err := json.Unmarshal(data, m); err != nil {
		return fmt.Errorf("failed to unmarshal message: %w", err)
	}
	return nil
}

// ToJSON converts a response to JSON bytes
func (r *Response) ToJSON() ([]byte, error) {
	data, err := json.Marshal(r)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal response: %w", err)
	}
	// Add newline delimiter for easy reading
	return append(data, '\n'), nil
}

// FromJSON parses a response from JSON bytes
func (r *Response) FromJSON(data []byte) error {
	if err := json.Unmarshal(data, r); err != nil {
		return fmt.Errorf("failed to unmarshal response: %w", err)
	}
	return nil
}

// Validate checks if a message is valid
func (m *Message) Validate() error {
	if m.Command == "" {
		return fmt.Errorf("command cannot be empty")
	}

	// Validate known commands
	validCommands := map[string]bool{
		CmdEcho:         true,
		CmdRegister:     true,
		CmdMessage:      true,
		CmdListUsers:    true,
		CmdListMessages: true,
		CmdTime:         true,
		CmdQuit:         true,
	}

	if !validCommands[m.Command] {
		return fmt.Errorf("unknown command: %s", m.Command)
	}

	// Some commands require data
	requiresData := map[string]bool{
		CmdEcho:     true,
		CmdRegister: true,
		CmdMessage:  true,
	}

	if requiresData[m.Command] && m.Data == "" {
		return fmt.Errorf("command %s requires data", m.Command)
	}

	return nil
}

// String returns a string representation of the message (for debugging)
func (m *Message) String() string {
	return fmt.Sprintf("Message{Command: %s, Data: %s, From: %s}", m.Command, m.Data, m.From)
}

// String returns a string representation of the response (for debugging)
func (r *Response) String() string {
	return fmt.Sprintf("Response{Success: %t, Message: %s, Data: %s}", r.Success, r.Message, r.Data)
}

package protocol

import (
	"testing"
)

func TestNewMessage(t *testing.T) {
	msg := NewMessage("ECHO", "Hello")

	if msg.Command != "ECHO" {
		t.Errorf("Expected command to be ECHO, got %s", msg.Command)
	}

	if msg.Data != "Hello" {
		t.Errorf("Expected data to be Hello, got %s", msg.Data)
	}
}

func TestMessageValidation(t *testing.T) {
	tests := []struct {
		name    string
		msg     Message
		wantErr bool
	}{
		{
			name:    "Valid ECHO message",
			msg:     Message{Command: "ECHO", Data: "test"},
			wantErr: false,
		},
		{
			name:    "Valid REGISTER message",
			msg:     Message{Command: "REGISTER", Data: "username"},
			wantErr: false,
		},
		{
			name:    "Empty command",
			msg:     Message{Command: "", Data: "test"},
			wantErr: true,
		},
		{
			name:    "Unknown command",
			msg:     Message{Command: "INVALID", Data: "test"},
			wantErr: true,
		},
		{
			name:    "ECHO without data",
			msg:     Message{Command: "ECHO", Data: ""},
			wantErr: true,
		},
		{
			name:    "QUIT is valid without data",
			msg:     Message{Command: "QUIT", Data: ""},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMessageJSON(t *testing.T) {
	msg := NewMessage("ECHO", "Hello World")

	// Test ToJSON
	jsonData, err := msg.ToJSON()
	if err != nil {
		t.Fatalf("ToJSON() error = %v", err)
	}

	// Test FromJSON
	var decoded Message
	err = decoded.FromJSON(jsonData)
	if err != nil {
		t.Fatalf("FromJSON() error = %v", err)
	}

	// Verify decoded message matches original
	if decoded.Command != msg.Command {
		t.Errorf("Command mismatch: got %s, want %s", decoded.Command, msg.Command)
	}

	if decoded.Data != msg.Data {
		t.Errorf("Data mismatch: got %s, want %s", decoded.Data, msg.Data)
	}
}

func TestResponseJSON(t *testing.T) {
	resp := NewResponse(true, "Success", "test data")

	// Test ToJSON
	jsonData, err := resp.ToJSON()
	if err != nil {
		t.Fatalf("ToJSON() error = %v", err)
	}

	// Test FromJSON
	var decoded Response
	err = decoded.FromJSON(jsonData)
	if err != nil {
		t.Fatalf("FromJSON() error = %v", err)
	}

	// Verify decoded response matches original
	if decoded.Success != resp.Success {
		t.Errorf("Success mismatch: got %v, want %v", decoded.Success, resp.Success)
	}

	if decoded.Message != resp.Message {
		t.Errorf("Message mismatch: got %s, want %s", decoded.Message, resp.Message)
	}

	if decoded.Data != resp.Data {
		t.Errorf("Data mismatch: got %s, want %s", decoded.Data, resp.Data)
	}
}

func TestMessageString(t *testing.T) {
	msg := Message{
		Command: "ECHO",
		Data:    "test",
		From:    "alice",
	}

	str := msg.String()
	expected := "Message{Command: ECHO, Data: test, From: alice}"

	if str != expected {
		t.Errorf("String() = %s, want %s", str, expected)
	}
}

func TestResponseString(t *testing.T) {
	resp := Response{
		Success: true,
		Message: "OK",
		Data:    "test",
	}

	str := resp.String()
	expected := "Response{Success: true, Message: OK, Data: test}"

	if str != expected {
		t.Errorf("String() = %s, want %s", str, expected)
	}
}

func BenchmarkMessageToJSON(b *testing.B) {
	msg := NewMessage("ECHO", "Benchmark test data")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = msg.ToJSON()
	}
}

func BenchmarkMessageFromJSON(b *testing.B) {
	msg := NewMessage("ECHO", "Benchmark test data")
	jsonData, _ := msg.ToJSON()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var decoded Message
		_ = decoded.FromJSON(jsonData)
	}
}

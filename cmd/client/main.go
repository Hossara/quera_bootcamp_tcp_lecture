package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"tcp_server/client"
	"tcp_server/protocol"
)

func main() {
	// Set up logging
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	// Server address
	address := "localhost:8080"
	if len(os.Args) > 1 {
		address = os.Args[1]
	}

	// Create client
	c := client.NewClient(address)

	// Connect to server
	fmt.Printf("🔌 Connecting to server at %s...\n", address)
	if err := c.Connect(); err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer c.Close()

	fmt.Printf("✅ Connected to server at %s\n\n", address)

	// Interactive loop
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Println("\n📋 Available commands:")
		fmt.Println("  1. REGISTER      - Register your username")
		fmt.Println("  2. MESSAGE       - Send a chat message")
		fmt.Println("  3. LIST_USERS    - List online users")
		fmt.Println("  4. ECHO          - Test echo")
		fmt.Println("  5. TIME          - Get server time")
		fmt.Println("  6. LIST_MESSAGES - List recent messages")
		fmt.Println("  7. QUIT          - Disconnect")
		fmt.Print("\nEnter command (or number): ")

		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}

		// Handle numeric shortcuts
		switch input {
		case "1":
			input = "REGISTER"
		case "2":
			input = "MESSAGE"
		case "3":
			input = "LIST_USERS"
		case "4":
			input = "ECHO"
		case "5":
			input = "TIME"
		case "6":
			input = "LIST_MESSAGES"
		case "7":
			input = "QUIT"
		}

		command := strings.ToUpper(input)

		// Handle commands that require data
		var data string
		switch command {
		case protocol.CmdRegister, protocol.CmdMessage, protocol.CmdEcho:
			fmt.Print("Enter data: ")
			if !scanner.Scan() {
				break
			}
			data = strings.TrimSpace(scanner.Text())
		}

		// Execute command
		switch command {
		case protocol.CmdRegister:
			if err := c.Register(data); err != nil {
				fmt.Printf("❌ Error: %v\n", err)
			}

		case protocol.CmdMessage:
			if err := c.SendChatMessage(data); err != nil {
				fmt.Printf("❌ Error: %v\n", err)
			}

		case protocol.CmdListUsers:
			if err := c.ListUsers(); err != nil {
				fmt.Printf("❌ Error: %v\n", err)
			}

		case protocol.CmdListMessages:
			if err := c.ListMessages(); err != nil {
				fmt.Printf("❌ Error: %v\n", err)
			}

		case protocol.CmdEcho:
			if err := c.Echo(data); err != nil {
				fmt.Printf("❌ Error: %v\n", err)
			}

		case protocol.CmdTime:
			if err := c.GetServerTime(); err != nil {
				fmt.Printf("❌ Error: %v\n", err)
			}

		case protocol.CmdQuit:
			fmt.Println("\n👋 Disconnecting...")
			c.Quit()
			return

		default:
			fmt.Printf("❌ Unknown command: %s\n", command)
		}
	}
}

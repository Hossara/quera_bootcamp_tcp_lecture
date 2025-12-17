package main

import (
	"fmt"
	"log"
	"tcp_server/client"
	"tcp_server/server"
	"time"
)

// This is a demonstration of TCP client-server communication
// For actual usage, run cmd/server/main.go and cmd/client/main.go separately
func main() {
	fmt.Println("🎓 TCP/IP Demo")
	fmt.Println("===========================\n")

	// Start server in background
	srv := server.NewServer(":9999")
	go func() {
		if err := srv.Start(); err != nil {
			log.Printf("Server error: %v", err)
		}
	}()

	// Give server time to start
	time.Sleep(1 * time.Second)

	// Create two clients
	fmt.Println("📱 Creating clients...")
	client1 := client.NewClient("localhost:9999")
	client2 := client.NewClient("localhost:9999")

	// Connect clients
	fmt.Println("\n🔌 Connecting Client 1...")
	if err := client1.Connect(); err != nil {
		log.Fatalf("Client 1 connection failed: %v", err)
	}
	defer func() { _ = client1.Close() }()

	fmt.Println("🔌 Connecting Client 2...")
	if err := client2.Connect(); err != nil {
		log.Fatalf("Client 2 connection failed: %v", err)
	}
	defer func() { _ = client2.Close() }()

	// Demo 1: Echo
	fmt.Println("\n\n📢 Demo 1: Echo Command")
	fmt.Println("------------------------")
	_ = client1.Echo("Hello from Client 1!")

	time.Sleep(500 * time.Millisecond)

	// Demo 2: Registration
	fmt.Println("\n\n✏️  Demo 2: User Registration")
	fmt.Println("------------------------------")
	_ = client1.Register("Alice")
	time.Sleep(500 * time.Millisecond)
	_ = client2.Register("Bob")

	time.Sleep(500 * time.Millisecond)

	// Demo 3: List Users
	fmt.Println("\n\n👥 Demo 3: List Online Users")
	fmt.Println("-----------------------------")
	_ = client1.ListUsers()

	time.Sleep(500 * time.Millisecond)

	// Demo 4: Get Server Time
	fmt.Println("\n\n🕐 Demo 4: Server Time")
	fmt.Println("-----------------------")
	_ = client1.GetServerTime()

	time.Sleep(500 * time.Millisecond)

	// Demo 5: Broadcast Messages
	fmt.Println("\n\n💬 Demo 5: Chat Messages (Broadcasting)")
	fmt.Println("----------------------------------------")
	_ = client1.SendChatMessage("Hi everyone! This is Alice.")
	time.Sleep(500 * time.Millisecond)
	_ = client2.SendChatMessage("Hey Alice! Bob here.")

	time.Sleep(1 * time.Second)

	// Cleanup
	fmt.Println("\n\n🛑 Shutting down...")
	_ = client1.Quit()
	_ = client2.Quit()
	srv.Shutdown()

	fmt.Println("\n✅ Demo completed!")
	fmt.Println("\n💡 To run interactive mode:")
	fmt.Println("   Terminal 1: go run cmd/server/main.go")
	fmt.Println("   Terminal 2: go run cmd/client/main.go")
	fmt.Println("   Terminal 3: go run cmd/client/main.go")
}

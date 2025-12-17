package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"tcp_server/server"
)

func main() {
	// Set up logging
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	// Server address
	address := ":8080"
	if len(os.Args) > 1 {
		address = os.Args[1]
	}

	// Create server
	srv := server.NewServer(address)
	defer srv.Shutdown()

	// Set up graceful shutdown on Ctrl+C
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Start server in goroutine
	go func() {
		if err := srv.Start(); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Wait for shutdown signal
	<-sigChan
	fmt.Println("\n\n🛑 Received shutdown signal...")
}

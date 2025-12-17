# TCP/IP Socket Programming in Go - Educational Repository

This repository is designed to teach the fundamentals of networking and socket programming using Go.

## 📚 Table of Contents

1. [How the Internet Works](#how-the-internet-works)
2. [What is TCP/IP](#what-is-tcpip)
3. [TCP vs UDP](#tcp-vs-udp)
4. [TCP Server Implementation](#tcp-server-implementation)
5. [TCP Client Implementation](#tcp-client-implementation)
6. [Real-World Example: Chat Server](#real-world-example-chat-server)
7. [Running the Examples](#running-the-examples)

---

## 🌐 How the Internet Works

The Internet is a global network of interconnected computers that communicate using standardized protocols.

### Key Concepts:

1. **Client-Server Model**: 
   - **Client**: Initiates requests (e.g., web browser, mobile app)
   - **Server**: Responds to requests (e.g., web server, database server)

2. **Network Layers** (OSI Model):
   - **Application Layer**: HTTP, FTP, SMTP, etc.
   - **Transport Layer**: TCP, UDP
   - **Network Layer**: IP (Internet Protocol)
   - **Data Link Layer**: Ethernet, WiFi
   - **Physical Layer**: Cables, wireless signals

3. **Data Flow**:
   ```
   Application (Browser) → TCP → IP → Network Interface → Physical Medium
   ↓
   Internet (Routers)
   ↓
   Physical Medium → Network Interface → IP → TCP → Application (Server)
   ```

4. **IP Addresses**: Unique identifiers for devices (e.g., 192.168.1.1, 2001:db8::1)

5. **Ports**: Endpoints for communication (e.g., HTTP uses port 80, HTTPS uses port 443)

---

## 🔌 What is TCP/IP

**TCP/IP** (Transmission Control Protocol/Internet Protocol) is the fundamental protocol suite that powers the Internet.

### IP (Internet Protocol):
- **Purpose**: Routes packets from source to destination
- **Addressing**: Uses IP addresses to identify devices
- **Connectionless**: Each packet is independent
- **Unreliable**: No guarantee of delivery

### TCP (Transmission Control Protocol):
- **Purpose**: Ensures reliable, ordered delivery of data
- **Connection-oriented**: Establishes a connection before data transfer
- **Features**:
  - **Reliability**: Guarantees data delivery
  - **Ordering**: Data arrives in the correct sequence
  - **Error Checking**: Detects and corrects errors
  - **Flow Control**: Prevents overwhelming the receiver
  - **Congestion Control**: Adapts to network conditions

### TCP Three-Way Handshake:
```
Client                    Server
  |                         |
  |--------- SYN --------->|  (Client initiates connection)
  |                         |
  |<----- SYN-ACK ---------|  (Server acknowledges)
  |                         |
  |--------- ACK --------->|  (Client confirms)
  |                         |
  |   Connection Established|
```

---

## ⚡ TCP vs UDP

| Feature | TCP | UDP |
|---------|-----|-----|
| **Connection** | Connection-oriented | Connectionless |
| **Reliability** | Guaranteed delivery | No guarantee |
| **Ordering** | Maintains order | No ordering |
| **Speed** | Slower (due to overhead) | Faster |
| **Header Size** | 20 bytes | 8 bytes |
| **Error Checking** | Extensive | Basic checksum |
| **Flow Control** | Yes | No |
| **Use Cases** | Web browsing, email, file transfer | Video streaming, online gaming, DNS |

### When to Use TCP:
- ✅ Data integrity is critical
- ✅ Order matters
- ✅ Connection state is needed
- Examples: HTTP, HTTPS, FTP, SSH, email

### When to Use UDP:
- ✅ Speed is more important than reliability
- ✅ Some data loss is acceptable
- ✅ Real-time communication
- Examples: Video streaming, VoIP, online gaming, DNS

---

## 🖥️ TCP Server Implementation

This repository demonstrates a TCP server that:
- Listens for incoming connections
- Handles multiple clients concurrently
- Receives structured messages (JSON)
- Processes requests and sends responses

### Key Components:

1. **Server Structure** (`server/server.go`):
   - Server configuration
   - Connection management
   - Message handling

2. **Protocol** (`protocol/protocol.go`):
   - Structured message format
   - Request/Response types
   - JSON encoding/decoding

3. **Concurrent Handling**:
   - Each client runs in a separate goroutine
   - Non-blocking I/O
   - Graceful connection cleanup

### Example:
```go
// Create and start server
srv := server.NewServer(":8080")
srv.Start()
```

---

## 💻 TCP Client Implementation

The TCP client demonstrates:
- Establishing TCP connections
- Sending structured messages
- Receiving and parsing responses
- Connection lifecycle management

### Key Features:

1. **Client Structure** (`client/client.go`):
   - Connection management
   - Message sending/receiving
   - Error handling

2. **Protocol Compliance**:
   - Same message format as server
   - JSON serialization
   - Type-safe communication

### Example:
```go
// Create client and connect
client := client.NewClient("localhost:8080")
client.Connect()

// Send message
response := client.SendMessage("ECHO", "Hello, Server!")
```

---

## 🚀 Real-World Example: Chat Server

This repository includes a complete chat server implementation demonstrating:

### Features:
- **Multi-client support**: Multiple users can connect simultaneously
- **Message broadcasting**: Messages sent to all connected clients
- **User management**: Username registration and tracking
- **Message history**: Server stores recent messages
- **Command system**: 
  - `REGISTER`: Set username
  - `MESSAGE`: Send chat message
  - `LIST_USERS`: Get list of online users
  - `LIST_MESSAGES`: Get recent message history
  - `ECHO`: Simple echo test
  - `TIME`: Get server time

### Architecture:
```
Server (Port 8080)
  ├── Client 1 (Goroutine)
  ├── Client 2 (Goroutine)
  ├── Client 3 (Goroutine)
  └── Message Broker (Broadcasts messages)
```

---

## 🏃 Running the Examples

### Prerequisites:
```bash
# Ensure Go is installed
go version

# Should output: go version go1.21 or higher
```

### 1. Start the Server:
```bash
# Run server from cmd/server
go run cmd/server/main.go

# Or build and run
go build -o server cmd/server/main.go
./server
```

### 2. Run Clients:

**Option A: Interactive Client**
```bash
go run cmd/client/main.go
```

**Option B: Multiple Clients**
Open multiple terminals and run the client in each:
```bash
# Terminal 1
go run cmd/client/main.go

# Terminal 2
go run cmd/client/main.go

# Terminal 3
go run cmd/client/main.go
```

### 3. Example Interaction:

**Client 1:**
```
Connected to server at localhost:8080
Enter command (REGISTER/MESSAGE/LIST_USERS/ECHO/TIME/QUIT): REGISTER
Enter data: Alice
Response: Registration successful. Welcome, Alice!

Enter command: MESSAGE
Enter data: Hello everyone!
Response: Message broadcasted
```

**Client 2:**
```
Connected to server at localhost:8080
Enter command: REGISTER
Enter data: Bob
Response: Registration successful. Welcome, Bob!

[Broadcast] Alice: Hello everyone!

Enter command: LIST_USERS
Response: Online users: Alice, Bob

Enter command: LIST_MESSAGES
Response: Recent messages:
[10:30:15] Alice: Hello everyone!
```

---

## 📂 Project Structure

```
tcp/
├── README.md                 # This file
├── CONCEPTS.md              # Detailed networking concepts
├── go.mod                   # Go module definition
├── main.go                  # Quick demo
├── protocol/
│   └── protocol.go          # Message protocol definition
├── server/
│   └── server.go            # TCP server implementation
├── client/
│   └── client.go            # TCP client implementation
└── cmd/
    ├── server/
    │   └── main.go          # Server entry point
    └── client/
        └── main.go          # Client entry point
```

---

## 🎓 Learning Path

1. **Start with concepts**: Read this README and CONCEPTS.md
2. **Examine protocol**: Understand the message format in `protocol/protocol.go`
3. **Study server**: Review `server/server.go` for server implementation
4. **Study client**: Review `client/client.go` for client implementation
5. **Run examples**: Start server and clients to see it in action
6. **Experiment**: Modify code, add features, break things and fix them!

---

## 🔧 Exercises

1. **Add authentication**: Require password for user registration
2. **Private messages**: Implement direct messaging between users
3. **Room system**: Add chat rooms/channels
4. **File transfer**: Implement file sending capability
5. **Encryption**: Add TLS for secure communication
6. **Protocol buffer**: Replace JSON with Protocol Buffers
7. **UDP comparison**: Implement same chat using UDP

---

## 📖 Additional Resources

- [Go net package documentation](https://pkg.go.dev/net)
- [TCP RFC 793](https://tools.ietf.org/html/rfc793)
- [IP RFC 791](https://tools.ietf.org/html/rfc791)
- [Beej's Guide to Network Programming](https://beej.us/guide/bgnet/)

### 📚 Documentation in This Repository

- **CONCEPTS.md** - Deep dive into networking concepts with detailed TCP/IP explanations
- **SEQUENCE_NUMBERS.md** - Quick reference guide for TCP sequence numbers
- **SEQUENCE_NUMBERS_VISUAL.md** - Visual diagrams and examples of sequence numbers
- **LIST_MESSAGES.md** - Documentation for the message history feature
- **EXERCISES.md** - Practice exercises and challenges
- **REAL_WORLD_EXAMPLES.md** - Real-world socket programming applications
- **QUICKSTART.md** - Get started in 5 minutes

---

## 📝 License

This repository is for educational purposes. Feel free to use, modify, and share.

---

## 🤝 Contributing

This is an educational project. Feel free to:
- Report issues
- Suggest improvements
- Add more examples
- Improve documentation

Happy Learning! 🎉


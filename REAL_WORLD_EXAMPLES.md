# Real-World Socket Programming Examples

This document provides real-world examples and use cases for socket programming.

---

## 🌐 1. Web Server (HTTP)

### Concept
Web servers use TCP sockets to serve HTTP requests. Every time you visit a website, your browser creates a TCP connection to the server.

### Simple HTTP Server Example

```go
package main

import (
	"bufio"
	"fmt"
	"net"
	"time"
)

func main() {
	listener, _ := net.Listen("tcp", ":8080")
	fmt.Println("HTTP Server listening on :8080")
	
	for {
		conn, _ := listener.Accept()
		go handleHTTP(conn)
	}
}

func handleHTTP(conn net.Conn) {
	defer conn.Close()
	
	// Read HTTP request
	reader := bufio.NewReader(conn)
	request, _ := reader.ReadString('\n')
	fmt.Printf("Request: %s", request)
	
	// Send HTTP response
	response := "HTTP/1.1 200 OK\r\n" +
		"Content-Type: text/html\r\n" +
		"\r\n" +
		"<html><body><h1>Hello from TCP Socket!</h1></body></html>"
	
	conn.Write([]byte(response))
}
```

**Test it**:
```bash
# Run the server
go run simple_http.go

# In another terminal
curl http://localhost:8080
# Or visit http://localhost:8080 in browser
```

---

## 💬 2. Real-Time Chat Application

### Architecture

```
┌─────────────┐     ┌──────────────┐     ┌─────────────┐
│  Client 1   │────▶│              │◀────│  Client 2   │
│  (Browser)  │     │ Chat Server  │     │  (Mobile)   │
└─────────────┘     │   (TCP)      │     └─────────────┘
                    │              │
┌─────────────┐     │  - Routing   │     ┌─────────────┐
│  Client 3   │────▶│  - Storage   │◀────│  Client 4   │
│  (Desktop)  │     │  - Auth      │     │  (Tablet)   │
└─────────────┘     └──────────────┘     └─────────────┘
```

### Features
- **Real-time messaging**: Messages delivered instantly
- **Presence**: Online/offline status
- **Typing indicators**: See when others are typing
- **Message history**: Retrieve past messages
- **Group chats**: Multiple participants

### Protocol Example

```json
// Message types
{
  "type": "message",
  "from": "alice",
  "to": "bob",
  "content": "Hello!",
  "timestamp": "2024-01-15T10:30:00Z"
}

{
  "type": "typing",
  "from": "alice",
  "to": "bob"
}

{
  "type": "presence",
  "user": "alice",
  "status": "online"
}
```

### Real Examples
- **Slack**: Uses WebSockets (upgrade from HTTP)
- **WhatsApp**: Custom protocol over TCP/UDP
- **Discord**: WebSockets for real-time, REST for API

---

## 🎮 3. Multiplayer Game Server

### Concept
Online games need fast, reliable communication between players and servers.

### Architecture Patterns

#### Pattern 1: Authoritative Server (Most Common)
```
┌─────────────┐
│  Player 1   │──┐
└─────────────┘  │
                 │    ┌──────────────────┐
┌─────────────┐  │    │  Game Server     │
│  Player 2   │──┼───▶│  - Game Logic    │
└─────────────┘  │    │  - State         │
                 │    │  - Validation    │
┌─────────────┐  │    └──────────────────┘
│  Player 3   │──┘
└─────────────┘

```

#### Pattern 2: Peer-to-Peer
```
┌─────────────┐
│  Player 1   │◀──────────────┐
└─────────────┘               │
       ▲                      │
       │                      │
       │         ┌─────────────▼──┐
       └─────────│  Player 2      │
                 └─────────────▲──┘
                               │
                               │
                 ┌─────────────┴──┐
                 │  Player 3      │
                 └────────────────┘
```

### Protocol Considerations

**TCP for**:
- Lobby/matchmaking
- Chat
- Inventory/transactions
- Critical events (player joined, game over)

**UDP for**:
- Player positions (fast, loss acceptable)
- Projectile positions
- Quick updates (30-60 times/second)

### Example: Simple Game State Sync

```go
type GameState struct {
	Players map[string]PlayerState
	Tick    int64
}

type PlayerState struct {
	X, Y     float64
	Health   int
	Rotation float64
}

// Client sends input
type Input struct {
	Tick      int64
	MoveX     float64  // -1 to 1
	MoveY     float64  // -1 to 1
	Action    string   // "shoot", "jump", etc.
}

// Server broadcasts state
func broadcastGameState(state GameState) {
	data, _ := json.Marshal(state)
	for _, player := range connectedPlayers {
		player.conn.Write(data)
	}
}
```

### Real Examples
- **Fortnite**: Uses UDP for gameplay, TCP for lobby
- **Minecraft**: TCP-based protocol
- **CS:GO**: Source Engine, mostly UDP

---

## 📊 4. Database Client-Server

### Concept
Databases like MySQL, PostgreSQL use custom TCP protocols.

### How It Works

```
┌──────────────┐                    ┌──────────────┐
│              │  1. Connect        │              │
│              │──────────────────▶│              │
│              │                    │              │
│   Client     │  2. Authenticate   │   Database   │
│ (Your App)   │──────────────────▶│   Server     │
│              │                    │              │
│              │  3. Query          │              │
│              │──────────────────▶│              │
│              │                    │              │
│              │  4. Results        │              │
│              │◀──────────────────│              │
└──────────────┘                    └──────────────┘
```

### PostgreSQL Wire Protocol Example

```
Client → Server: StartupMessage
Server → Client: AuthenticationRequest
Client → Server: PasswordMessage
Server → Client: AuthenticationOK
Server → Client: ReadyForQuery

Client → Server: Query("SELECT * FROM users")
Server → Client: RowDescription (column info)
Server → Client: DataRow (row 1)
Server → Client: DataRow (row 2)
...
Server → Client: CommandComplete
Server → Client: ReadyForQuery
```

### Connection Pooling

Instead of creating a new connection for each query:

```go
// Bad: Create connection every time
func getUser(id int) User {
	conn, _ := sql.Open("postgres", "...")
	defer conn.Close()
	// query...
}

// Good: Use connection pool
var db *sql.DB // Global pool

func init() {
	db, _ = sql.Open("postgres", "...")
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
}

func getUser(id int) User {
	// Uses connection from pool
	db.QueryRow("SELECT * FROM users WHERE id = $1", id)
}
```

---

## 📧 5. Email (SMTP/IMAP)

### SMTP (Sending Email)

```
Client → Server: HELO client.example.com
Server → Client: 250 Hello client.example.com

Client → Server: MAIL FROM:<sender@example.com>
Server → Client: 250 OK

Client → Server: RCPT TO:<recipient@example.com>
Server → Client: 250 OK

Client → Server: DATA
Server → Client: 354 Start mail input

Client → Server: Subject: Test Email
Client → Server: 
Client → Server: This is the email body.
Client → Server: .
Server → Client: 250 OK: Message accepted

Client → Server: QUIT
Server → Client: 221 Bye
```

### IMAP (Reading Email)

```
Client → Server: LOGIN username password
Server → Client: OK LOGIN completed

Client → Server: SELECT INBOX
Server → Client: * 42 EXISTS
Server → Client: OK SELECT completed

Client → Server: FETCH 1 BODY[]
Server → Client: * 1 FETCH (BODY[] {342}
                [email content]
                )
Server → Client: OK FETCH completed
```

---

## 🏭 6. IoT & Telemetry

### Concept
IoT devices (sensors, smart home) send data to servers.

### Architecture

```
┌────────────┐  ┌────────────┐  ┌────────────┐
│ Temperature│  │  Motion    │  │   Camera   │
│  Sensor    │  │  Sensor    │  │            │
└──────┬─────┘  └──────┬─────┘  └──────┬─────┘
       │                │                │
       │    TCP/MQTT    │                │
       └────────┬───────┴────────────────┘
                │
       ┌────────▼────────┐
       │  IoT Gateway    │
       │  - Aggregation  │
       │  - Buffering    │
       └────────┬────────┘
                │
                │ Internet
                │
       ┌────────▼────────┐
       │  Cloud Server   │
       │  - Storage      │
       │  - Analytics    │
       └─────────────────┘
```

### MQTT Protocol (Built on TCP)

```go
// MQTT is a lightweight pub/sub protocol

// Publisher (Sensor)
client.Publish("home/temperature", "22.5")
client.Publish("home/humidity", "65")

// Subscriber (Dashboard)
client.Subscribe("home/#") // # is wildcard

// Receives:
// home/temperature: 22.5
// home/humidity: 65
```

### Real Example: Smart Home

```go
type TelemetryData struct {
	DeviceID  string
	Timestamp time.Time
	Type      string // "temperature", "motion", etc.
	Value     float64
}

func handleSensor(conn net.Conn) {
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		var data TelemetryData
		json.Unmarshal(scanner.Bytes(), &data)
		
		// Store in database
		db.Insert(data)
		
		// Check thresholds
		if data.Type == "temperature" && data.Value > 30 {
			sendAlert("Temperature too high!")
		}
	}
}
```

---

## 🎥 7. Video Streaming

### Live Streaming Architecture

```
┌─────────────┐
│   Camera    │
└──────┬──────┘
       │ RTMP (TCP)
       │
┌──────▼──────┐
│   Encoder   │
└──────┬──────┘
       │
┌──────▼──────────────────┐
│  Streaming Server       │
│  - Transcoding          │
│  - Adaptive Bitrate     │
└──────┬──────────────────┘
       │
       ├─────────────────┬─────────────────┐
       │                 │                 │
┌──────▼──────┐  ┌──────▼──────┐  ┌──────▼──────┐
│  Viewer 1   │  │  Viewer 2   │  │  Viewer 3   │
│  (HLS/TCP)  │  │  (HLS/TCP)  │  │  (WebRTC)   │
└─────────────┘  └─────────────┘  └─────────────┘
```

### Protocols

**RTMP (Real-Time Messaging Protocol)** - TCP-based
- Used for ingesting stream
- Reliable delivery
- Higher latency (3-5 seconds)

**HLS (HTTP Live Streaming)** - HTTP/TCP
- Chunks video into small segments
- Served over HTTP
- Adaptive bitrate
- Latency: 10-30 seconds

**WebRTC** - UDP-based
- Peer-to-peer
- Very low latency (<1 second)
- Used for video calls

---

## 🔒 8. VPN (Virtual Private Network)

### Concept
VPN creates an encrypted tunnel between client and server.

```
┌──────────────┐                           ┌──────────────┐
│   Your PC    │                           │  VPN Server  │
│              │   Encrypted TCP Tunnel    │              │
│  ┌────────┐  │═══════════════════════════│  ┌────────┐  │
│  │ App    │──┼───────────────────────────┼─▶│Internet│  │
│  └────────┘  │                           │  └────────┘  │
└──────────────┘                           └──────────────┘
```

### How it Works

1. **Establish TCP connection** to VPN server
2. **TLS handshake** for encryption
3. **All traffic** routes through this tunnel
4. **Server forwards** to actual destination

---

## 🔐 9. SSH (Secure Shell)

### Concept
Remote access to servers over encrypted channel.

### Protocol Flow

```
Client → Server: SSH Version Exchange
Client ↔ Server: Key Exchange (Diffie-Hellman)
Client ↔ Server: Algorithm Negotiation
Client → Server: Authentication (password/key)
Server → Client: Authentication Success

Client → Server: Request Shell
Client ↔ Server: Interactive Session (encrypted)
```

### Example Usage

```bash
# SSH creates a TCP connection
ssh user@server.com

# Behind the scenes:
# 1. TCP connection to port 22
# 2. Encryption handshake
# 3. Authentication
# 4. Encrypted channel for all commands
```

---

## 📱 10. Mobile Push Notifications

### Architecture

```
┌──────────────┐
│  Your Server │
└──────┬───────┘
       │ HTTPS
       │
┌──────▼──────────────────┐
│  Push Service           │
│  - Apple APNs (TCP)     │
│  - Google FCM (HTTP/2)  │
└──────┬──────────────────┘
       │
       ├─────────────┬─────────────┐
       │             │             │
┌──────▼──────┐ ┌───▼──────┐ ┌────▼──────┐
│   iPhone    │ │  Android │ │  Android  │
└─────────────┘ └──────────┘ └───────────┘
```

### How It Works

1. **Device registers** with push service
2. **Gets device token**
3. **App sends token** to your server
4. **Your server sends** notification to push service
5. **Push service delivers** to device (maintains persistent TCP connection)

---

## 🛠️ Protocol Comparison Summary

| Use Case | Protocol | Reason |
|----------|----------|--------|
| Web browsing | HTTP/TCP | Reliability needed |
| Live video | WebRTC/UDP | Speed > perfect quality |
| Video streaming | HLS/TCP | Can buffer, needs quality |
| Online gaming | UDP + TCP | Both - UDP for position, TCP for critical |
| Chat | WebSocket/TCP | Real-time + reliability |
| Email | SMTP/TCP | Must deliver all content |
| File transfer | FTP/TCP | Complete file needed |
| IoT telemetry | MQTT/TCP | Reliable data delivery |
| Video calls | WebRTC/UDP | Real-time is critical |
| Database | Custom/TCP | Reliability essential |

---

## 🎓 Key Takeaways

1. **Choose the right protocol** for your use case
2. **TCP for reliability**, UDP for speed
3. **Consider latency** requirements
4. **Design your application protocol** carefully
5. **Handle errors gracefully**
6. **Scale with connection pooling**, load balancing
7. **Security matters**: Use TLS/SSL
8. **Monitor and measure** performance

---

## 🔍 Further Exploration

### Tools to Explore

1. **Wireshark**: Packet analyzer
2. **netcat**: TCP/UDP testing
3. **tcpdump**: Command-line packet capture
4. **curl**: HTTP client testing
5. **telnet**: Raw TCP connection testing

### Try It Yourself

```bash
# Connect to a web server manually
telnet example.com 80

# Then type:
GET / HTTP/1.1
Host: example.com

# Press Enter twice

# You'll see the raw HTTP response!
```

---

Happy learning! 🚀 These examples show how TCP/IP powers the modern Internet. Your chat server project is built on the same principles used by companies like Facebook, Google, and Amazon!


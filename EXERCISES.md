# Exercises and Challenges

## 🎯 Beginner Exercises

### 1. Add a PING Command
**Goal**: Implement a simple PING/PONG mechanism

**Tasks**:
- Add `CmdPing = "PING"` to protocol constants
- Server should respond with "PONG" and current timestamp
- Client should calculate round-trip time

**Hint**: Use `time.Now()` before sending and after receiving

---

### 2. Username Validation
**Goal**: Add validation for usernames

**Tasks**:
- Username must be 3-20 characters
- Only alphanumeric characters allowed
- No duplicate usernames

**Files to modify**: `server/server.go` in the `CmdRegister` case

---

### 3. Command History
**Goal**: Keep track of commands sent by each client

**Tasks**:
- Add a `history []string` field to the Client struct
- Store last 10 commands
- Add `LIST_HISTORY` command to view history

---

## 🚀 Intermediate Exercises

### 4. Private Messaging
**Goal**: Allow users to send messages to specific users

**Tasks**:
- Add `CmdPrivateMessage = "PM"` command
- Message format: `username:message`
- Parse recipient from message data
- Send only to specific user

**Example**:
```
Client 1: PM Bob:Hi Bob, this is private!
```

---

### 5. Chat Rooms/Channels
**Goal**: Implement multiple chat rooms

**Tasks**:
- Add `CmdJoinRoom = "JOIN_ROOM"` command
- Add `CmdLeaveRoom = "LEAVE_ROOM"` command
- Modify broadcast to send only to users in the same room
- Add `LIST_ROOMS` command

**Data structures needed**:
```go
type Room struct {
    name    string
    clients map[*Client]bool
}
```

---

### 6. Message Persistence
**Goal**: Save chat history to a file

**Tasks**:
- Create a `messages.log` file
- Append all messages to the file
- Add `GET_HISTORY` command to retrieve last N messages
- Use `os.OpenFile()` with `os.O_APPEND`

---

### 7. Authentication System
**Goal**: Require password for registration

**Tasks**:
- Modify `REGISTER` to accept `username:password`
- Hash passwords using `crypto/sha256`
- Add `LOGIN` command for returning users
- Store credentials in a map or file

---

## 🔥 Advanced Exercises

### 8. TLS/SSL Encryption
**Goal**: Secure the connection with TLS

**Tasks**:
- Generate self-signed certificates
- Use `net.ListenTLS()` instead of `net.Listen()`
- Use `tls.Dial()` instead of `net.Dial()`
- Update client and server

**Resources**: `crypto/tls` package

---

### 9. File Transfer
**Goal**: Allow users to send files to each other

**Tasks**:
- Add `SEND_FILE` command
- Implement chunking for large files
- Add progress indicator
- Handle concurrent file transfers

**Protocol considerations**:
- File metadata (name, size)
- Chunk numbering
- Error handling and retransmission

---

### 10. Rate Limiting
**Goal**: Prevent spam by limiting message rate

**Tasks**:
- Implement token bucket algorithm
- Limit to 5 messages per minute per user
- Return error if limit exceeded
- Add admin bypass

**Data structure**:
```go
type RateLimiter struct {
    tokens     int
    lastRefill time.Time
}
```

---

### 11. Admin Commands
**Goal**: Add administrative capabilities

**Tasks**:
- Add admin authentication
- Implement `KICK` command (disconnect user)
- Implement `BAN` command (prevent reconnection)
- Implement `BROADCAST_ALL` (server-wide announcement)
- Add `SERVER_STATS` (connection count, uptime, etc.)

---

### 12. Protocol Buffers
**Goal**: Replace JSON with Protocol Buffers for efficiency

**Tasks**:
- Install protoc compiler
- Define `.proto` file for messages
- Generate Go code
- Update client and server to use protobuf
- Compare performance with JSON

**Example .proto**:
```protobuf
syntax = "proto3";

message Message {
    string command = 1;
    string data = 2;
    string from = 3;
}
```

---

## 🏆 Expert Challenges

### 13. Load Balancer
**Goal**: Create a load balancer for multiple servers

**Tasks**:
- Run multiple server instances
- Implement load balancer that distributes clients
- Use round-robin or least-connections algorithm
- Handle server failures

---

### 14. Distributed Chat (Pub/Sub)
**Goal**: Allow multiple servers to communicate

**Tasks**:
- Integrate Redis Pub/Sub
- Each server subscribes to channels
- Messages published to Redis are received by all servers
- Broadcast to all connected clients across servers

---

### 15. WebSocket Bridge
**Goal**: Allow web browsers to connect

**Tasks**:
- Create WebSocket server (port 8081)
- Bridge WebSocket messages to TCP server
- Create simple HTML/JavaScript client
- Handle protocol translation

---

### 16. Custom Binary Protocol
**Goal**: Design efficient binary protocol

**Tasks**:
- Design binary message format
- Implement encoding/decoding
- Use fixed-size headers
- Compare with JSON (size, speed)

**Example format**:
```
[1 byte: version][1 byte: command][4 bytes: length][N bytes: data]
```

---

### 17. Connection Pooling
**Goal**: Implement client-side connection pool

**Tasks**:
- Create pool with min/max connections
- Reuse idle connections
- Handle connection lifecycle
- Add timeout for idle connections

---

### 18. Metrics and Monitoring
**Goal**: Add observability

**Tasks**:
- Instrument code with metrics
- Track: connections, messages/sec, errors
- Expose metrics endpoint (HTTP)
- Create Prometheus exporter
- Build Grafana dashboard

---

## 📊 Testing Exercises

### 19. Unit Tests
**Goal**: Add comprehensive test coverage

**Tasks**:
- Write tests for protocol encoding/decoding
- Test message validation
- Mock network connections
- Aim for 80%+ coverage

---

### 20. Load Testing
**Goal**: Test server performance

**Tasks**:
- Create load test client
- Simulate 1000+ concurrent connections
- Measure latency and throughput
- Identify bottlenecks
- Use `pprof` for profiling

**Tools**: `hey`, `wrk`, or custom Go tool

---

## 🎓 Learning Projects

### 21. Build a Complete Application

Choose one:

**A. Multiplayer Game**
- Simple turn-based game (tic-tac-toe, chess)
- Real-time game state synchronization
- Matchmaking system

**B. Collaborative Text Editor**
- Multiple users editing same document
- Conflict resolution
- Operational transformation or CRDTs

**C. IoT Device Simulator**
- Simulate sensor devices
- Send telemetry data
- Command and control
- Data aggregation server

**D. Distributed Task Queue**
- Clients submit tasks
- Workers pull and execute tasks
- Result collection
- Priority queue

---

## 💡 Debugging Exercises

### 22. Packet Analysis
**Goal**: Understand TCP at packet level

**Tasks**:
- Use Wireshark to capture packets
- Identify TCP handshake
- See your messages in packets
- Analyze retransmissions and timeouts

---

### 23. Error Injection
**Goal**: Test error handling

**Tasks**:
- Simulate network errors
- Test connection drops
- Test malformed messages
- Test timeout scenarios
- Ensure graceful degradation

---

## 📝 Submission Guidelines

For each exercise:

1. **Code**: Implement the feature
2. **Tests**: Write tests
3. **Documentation**: Explain your approach
4. **Demo**: Show it working

Example submission structure:
```
exercises/
  exercise-04-private-messaging/
    README.md
    protocol.go
    server.go
    client.go
    demo.go
    *_test.go
```

---

## 🎉 Bonus: Real-World Integration

### 24. Cloud Deployment
- Deploy to AWS/GCP/Azure
- Use Docker containers
- Set up CI/CD pipeline
- Configure load balancing

### 25. Database Integration
- Store users in PostgreSQL
- Save message history
- Implement user profiles
- Add search functionality

---

## 📚 Resources

- [Go net package](https://pkg.go.dev/net)
- [Go concurrency patterns](https://go.dev/blog/pipelines)
- [TCP RFC](https://tools.ietf.org/html/rfc793)
- [Protocol Buffers](https://developers.google.com/protocol-buffers)

Happy coding! 🚀


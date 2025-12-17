# Quick Start Guide

## 🚀 Getting Started in 5 Minutes

### Prerequisites
- Go 1.20 or higher installed
- Terminal/Command Prompt

### Installation

```bash
# Clone or download this repository
cd tcp

# Verify Go installation
go version
```

### Run the Quick Demo

```bash
# Run the automated demo
go run main.go
```

This will start a server and two clients automatically, demonstrating all features.

---

## 🎮 Interactive Mode

### Option 1: Terminal Mode (Recommended for Learning)

**Terminal 1 - Start Server:**
```bash
go run cmd/server/main.go
```

You should see:
```
🚀 TCP Server started on :8080
📖 Educational TCP/IP Server - Ready to accept connections
```

**Terminal 2 - Start First Client:**
```bash
go run cmd/client/main.go
```

**Terminal 3 - Start Second Client (optional):**
```bash
go run cmd/client/main.go
```

### Option 2: Build and Run Binaries

```bash
# Build server and client
go build -o server cmd/server/main.go
go build -o client cmd/client/main.go

# Run server
./server

# Run client (in another terminal)
./client
```

---

## 📱 Using the Client

When you run the client, you'll see this menu:

```
📋 Available commands:
  1. REGISTER   - Register your username
  2. MESSAGE    - Send a chat message
  3. LIST_USERS - List online users
  4. ECHO       - Test echo
  5. TIME       - Get server time
  6. QUIT       - Disconnect

Enter command (or number):
```

### Example Session

```
> 1                    # Choose REGISTER
Enter data: Alice
✅ Registration successful. Welcome, Alice!

> 2                    # Choose MESSAGE
Enter data: Hello everyone!
✅ Message broadcasted

> 3                    # Choose LIST_USERS
👥 Online users: Alice, Bob

> 6                    # Choose QUIT
👋 Goodbye!
```

---

## 🔧 Configuration

### Change Server Port

Edit `cmd/server/main.go`:
```go
address := ":8080"  // Change to ":9000" or any port
```

Or pass as argument:
```bash
go run cmd/server/main.go :9000
```

### Connect to Different Server

Edit `cmd/client/main.go`:
```go
address := "localhost:8080"  // Change to "server.com:8080"
```

Or pass as argument:
```bash
go run cmd/client/main.go server.com:8080
```

---

## 🧪 Testing Features

### Test Echo
```
> ECHO
Enter data: Hello Server!
📢 Echo: Hello Server!
```

### Test Broadcasting
1. Open 2+ clients
2. Register different usernames
3. Send messages from one client
4. See messages appear on other clients

### Test User List
```
> LIST_USERS
👥 Online users: Alice, Bob, Charlie
```

### Test Server Time
```
> TIME
🕐 Server time: 2024-01-15T10:30:00Z
```

---

## 🐛 Troubleshooting

### "connection refused"
- Make sure server is running first
- Check if port is already in use: `lsof -i :8080` (Linux/Mac) or `netstat -an | findstr 8080` (Windows)

### "address already in use"
- Another program is using the port
- Change the port number
- Kill the existing process

### Client hangs/freezes
- Press Ctrl+C to exit
- Check network connection
- Restart server and client

---

## 📖 Next Steps

1. **Read the docs**: Check `README.md` for detailed concepts
2. **Study the code**: 
   - Start with `protocol/protocol.go`
   - Then `server/server.go`
   - Then `client/client.go`
3. **Explore concepts**: Read `CONCEPTS.md`
4. **Try exercises**: See `EXERCISES.md`
5. **Real-world examples**: Read `REAL_WORLD_EXAMPLES.md`

---

## 🎓 Learning Path

### Day 1: Basics
- Run the demo
- Understand the protocol
- Read about TCP/IP

### Day 2: Server
- Study server implementation
- Understand concurrency (goroutines)
- Learn about connection handling

### Day 3: Client
- Study client implementation
- Understand protocol compliance
- Learn about error handling

### Day 4: Practice
- Try exercises
- Modify the code
- Add new features

### Day 5: Advanced
- Read about real-world examples
- Study performance optimization
- Explore security (TLS)

---

## 💡 Tips

1. **Use multiple terminals**: One for server, multiple for clients
2. **Read server logs**: They show what's happening
3. **Experiment**: Try breaking things and fixing them
4. **Use Wireshark**: See actual TCP packets
5. **Test edge cases**: What happens if you send invalid data?

---

## 📚 Resources in This Repo

- `README.md` - Main documentation
- `CONCEPTS.md` - Detailed networking concepts
- `EXERCISES.md` - Practice exercises
- `REAL_WORLD_EXAMPLES.md` - Real-world use cases
- `protocol/` - Message protocol
- `server/` - Server implementation
- `client/` - Client implementation
- `cmd/` - Command-line applications

---

## 🤝 Need Help?

1. Check the documentation files
2. Read the code comments
3. Try the exercises
4. Experiment and learn!

Happy Learning! 🎉


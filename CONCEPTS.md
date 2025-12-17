# Deep Dive: Networking Concepts

## 📡 How the Internet Really Works

### The Journey of a Web Request

When you type `www.example.com` in your browser:

1. **DNS Resolution**:
   - Browser asks DNS server: "What's the IP address of example.com?"
   - DNS responds: "It's 93.184.216.34"

2. **TCP Connection** (Three-Way Handshake):
   ```
   Your Computer (192.168.1.100:54321) → Server (93.184.216.34:80)
   
   Step 1: SYN (Synchronize)
   - "Hello server, I want to connect. My sequence number is 1000"
   
   Step 2: SYN-ACK (Synchronize-Acknowledge)
   - "Hi! I accept. My sequence number is 5000. I acknowledge your 1000"
   
   Step 3: ACK (Acknowledge)
   - "Great! I acknowledge your 5000. Let's talk!"
   ```

   **Understanding Sequence Numbers:**
   
   The **sequence number** is a 32-bit counter in the TCP header that tracks bytes transmitted. Each side picks an initial sequence number (ISN) when establishing a connection.
   
   **Why Sequence Numbers?**
   - **Ordering**: Ensures data arrives in correct order
   - **Reliability**: Detects missing or duplicated packets
   - **Flow Control**: Tracks what's been sent and acknowledged
   
   **How It Works:**
   
   Initial State:
   ```
   Client chooses ISN: 1000
   Server chooses ISN: 5000
   ```
   
   During Handshake:
   ```
   Client → Server: SYN, Seq=1000
   (SYN consumes 1 sequence number)
   
   Server → Client: SYN-ACK, Seq=5000, Ack=1001
   (Server expects client's next byte at 1001)
   (Server's SYN also consumes 1 sequence number)
   
   Client → Server: ACK, Seq=1001, Ack=5001
   (Both sides now ready to send data)
   ```
   
   **When Does Sequence Number Increase?**
   
   1. **SYN Flag**: Consumes 1 sequence number
   2. **FIN Flag**: Consumes 1 sequence number
   3. **Data Bytes**: Each byte of payload increases sequence number
   4. **Pure ACK**: Does NOT increase sequence number (no data)
   
   **Example with Data Transfer:**
   ```
   Client sends 100 bytes:
   Packet 1: Seq=1001, Data="Hello..." (100 bytes)
   
   Server acknowledges:
   Packet 2: Ack=1101 (expecting byte 1101 next)
   
   Client sends 50 more bytes:
   Packet 3: Seq=1101, Data="World..." (50 bytes)
   
   Server acknowledges:
   Packet 4: Ack=1151 (expecting byte 1151 next)
   ```
   
   **Bidirectional Communication:**
   ```
   Time    Client (Seq/Ack)              Server (Seq/Ack)
   ----    ----------------              ----------------
   T1      Seq=1001, 100 bytes    →
   T2                             ←      Ack=1101
   T3                             ←      Seq=5001, 200 bytes
   T4      Ack=5201               →
   T5      Seq=1101, 50 bytes     →
   T6                             ←      Ack=1151, Seq=5201, 150 bytes
   ```
   
   **Why Sequence Numbers Increase:**
   
   Think of sequence numbers as page numbers in a book:
   - They tell you which "page" (byte) of data you're looking at
   - If pages arrive out of order, you can reorder them
   - If a page is missing, you know which one to request again
   
   **Real Example:**
   ```
   Sending "HELLO WORLD" (11 bytes):
   
   Initial: Seq=1000 (after handshake: Seq=1001)
   
   Send "HELLO" (5 bytes):
   Packet: Seq=1001, Data="HELLO"
   Next Seq will be: 1001 + 5 = 1006
   
   Send " WORLD" (6 bytes):
   Packet: Seq=1006, Data=" WORLD"
   Next Seq will be: 1006 + 6 = 1012
   ```
   
   **Acknowledgment Number:**
   
   The ACK number tells the sender: "I've received everything up to this byte, send me byte X next"
   
   ```
   Client sends Seq=1001, 100 bytes
   Server responds with Ack=1101
   Meaning: "I got bytes 1001-1100, send 1101 next"
   ```
   
   **Lost Packet Example:**
   ```
   Client sends: Seq=1001, 100 bytes
   (Packet lost in network!)
   
   Server: (timeout, no ACK received)
   
   Client retransmits: Seq=1001, 100 bytes (same sequence!)
   
   Server receives: Ack=1101 (acknowledges receipt)
   ```
   
   **Out-of-Order Delivery:**
   ```
   Client sends:
   - Packet A: Seq=1001, 100 bytes (gets delayed)
   - Packet B: Seq=1101, 100 bytes (arrives first)
   
   Server receives Packet B first:
   - Stores it in buffer
   - Sends Ack=1001 (still waiting for Packet A)
   
   Server receives Packet A:
   - Combines with Packet B
   - Sends Ack=1201 (got everything up to 1200)
   ```
   
   **Key Insights:**
   
   1. **32-bit Number**: Can count up to 4,294,967,295 bytes (~4GB)
   2. **Wraps Around**: After reaching max, wraps back to 0
   3. **Each Direction Independent**: Client and server have separate sequence numbers
   4. **Random ISN**: Initial sequence number is randomized for security
   5. **Byte-Level Tracking**: Tracks individual bytes, not packets

3. **HTTP Request** (Application Layer):
   ```
   GET / HTTP/1.1
   Host: www.example.com
   ```

4. **Server Processing**:
   - Server receives request
   - Processes it (queries database, runs code, etc.)
   - Prepares response

5. **HTTP Response**:
   ```
   HTTP/1.1 200 OK
   Content-Type: text/html
   
   <html>...</html>
   ```

6. **TCP Connection Close**:
   ```
   Step 1: FIN (Finish)
   Step 2: ACK (Acknowledge)
   Step 3: FIN (Finish)
   Step 4: ACK (Acknowledge)
   ```

### Routing Through the Internet

Your data doesn't go directly to the destination. It hops through multiple routers:

```
Your Computer → Home Router → ISP Router → Internet Backbone → 
Destination ISP → Destination Server
```

Each hop uses IP to determine the next destination.

---

## 🔍 TCP/IP in Detail

### TCP Segment Structure

```
 0                   1                   2                   3
 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|          Source Port          |       Destination Port        |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|                        Sequence Number                        |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|                    Acknowledgment Number                      |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|  Data |           |U|A|P|R|S|F|                               |
| Offset| Reserved  |R|C|S|S|Y|I|            Window             |
|       |           |G|K|H|T|N|N|                               |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|           Checksum            |         Urgent Pointer        |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
```

**Key Fields**:
- **Source/Destination Port**: Identifies applications
- **Sequence Number**: Orders data
- **Acknowledgment Number**: Confirms received data
- **Flags**: SYN, ACK, FIN, etc.
- **Window**: Flow control
- **Checksum**: Error detection

### TCP Features Explained

#### 1. Reliability (Acknowledgment + Retransmission)

```
Sender                          Receiver
  |                               |
  |------- Data (Seq=100) ------->|
  |                               | ✓ Received
  |<------ ACK (Ack=101) ---------|
  |                               |
  |------- Data (Seq=101) ------->| ✗ Lost!
  |                               |
  |  (Timeout - No ACK received)  |
  |                               |
  |------- Data (Seq=101) ------->| (Retransmission)
  |                               | ✓ Received
  |<------ ACK (Ack=102) ---------|
```

#### 2. Flow Control (Sliding Window)

```
Sender has 10 packets to send
Receiver window size = 3 (can only handle 3 at a time)

Sender:  [1][2][3][4][5][6][7][8][9][10]
         Send →

Round 1: Send packets 1, 2, 3
Round 2: Wait for ACKs, then send 4, 5, 6
Round 3: Wait for ACKs, then send 7, 8, 9
Round 4: Wait for ACKs, then send 10
```

#### 3. Congestion Control

TCP adjusts sending rate based on network conditions:

```
Slow Start → Congestion Avoidance → Congestion Detected → Reduce Rate
     ↑                                                           |
     |___________________________________________________________|
                         (Recovery)
```

---

## 🔢 Deep Dive: TCP Sequence Numbers

### What Are Sequence Numbers?

TCP sequence numbers are 32-bit integers that label each byte of data transmitted in a TCP connection. They serve as unique identifiers for data bytes, enabling TCP to provide reliable, ordered delivery.

### The Anatomy of Sequence Numbers

```
TCP Header (simplified):
┌─────────────────────────────────────┐
│ Sequence Number (32 bits)          │ ← Identifies this packet's data
├─────────────────────────────────────┤
│ Acknowledgment Number (32 bits)    │ ← Confirms received data
└─────────────────────────────────────┘

Range: 0 to 4,294,967,295 (2^32 - 1)
```

### Initial Sequence Number (ISN)

When a TCP connection starts, each side randomly chooses an Initial Sequence Number:

```
Why Random?
1. Security: Prevents hijacking attacks
2. Uniqueness: Avoids confusion with old connections
3. Unpredictability: Makes connection harder to spoof

Example:
Client ISN: 1,458,234,567
Server ISN: 2,891,456,123
```

### Complete Example: From Handshake to Data Transfer

```
HANDSHAKE PHASE:
═══════════════

Step 1: Client → Server (SYN)
┌────────────────────────────┐
│ SYN Flag: 1                │
│ Seq: 1000                  │ ← Client's ISN
│ Ack: 0                     │
│ Data: (none)               │
└────────────────────────────┘
Effect: Seq 1000 is consumed by SYN flag
        Next available: 1001


Step 2: Server → Client (SYN-ACK)
┌────────────────────────────┐
│ SYN Flag: 1, ACK Flag: 1   │
│ Seq: 5000                  │ ← Server's ISN
│ Ack: 1001                  │ ← Acknowledges client's SYN (1000+1)
│ Data: (none)               │
└────────────────────────────┘
Effect: Seq 5000 is consumed by SYN flag
        Next available: 5001


Step 3: Client → Server (ACK)
┌────────────────────────────┐
│ ACK Flag: 1                │
│ Seq: 1001                  │ ← Client's next sequence
│ Ack: 5001                  │ ← Acknowledges server's SYN (5000+1)
│ Data: (none)               │
└────────────────────────────┘
Effect: Connection established
        Client ready at Seq 1001
        Server ready at Seq 5001


DATA TRANSFER PHASE:
═══════════════════

Transfer 1: Client sends "HELLO" (5 bytes)
┌────────────────────────────┐
│ Seq: 1001                  │ ← Starting byte number
│ Ack: 5001                  │
│ Data: "HELLO" (5 bytes)    │
└────────────────────────────┘
Effect: Bytes 1001-1005 sent
        Next available: 1006


Transfer 2: Server acknowledges
┌────────────────────────────┐
│ Seq: 5001                  │
│ Ack: 1006                  │ ← "I got up to 1005, send 1006 next"
│ Data: (none - pure ACK)    │
└────────────────────────────┘
Effect: Server Seq stays at 5001 (no data sent)


Transfer 3: Server sends "WORLD" (5 bytes)
┌────────────────────────────┐
│ Seq: 5001                  │ ← Starting byte number
│ Ack: 1006                  │
│ Data: "WORLD" (5 bytes)    │
└────────────────────────────┘
Effect: Bytes 5001-5005 sent
        Next available: 5006


Transfer 4: Client acknowledges
┌────────────────────────────┐
│ Seq: 1006                  │
│ Ack: 5006                  │ ← "I got up to 5005, send 5006 next"
│ Data: (none - pure ACK)    │
└────────────────────────────┘
Effect: Client Seq stays at 1006 (no data sent)


Transfer 5: Client sends " TCP!" (5 bytes)
┌────────────────────────────┐
│ Seq: 1006                  │ ← Continues from last data
│ Ack: 5006                  │
│ Data: " TCP!" (5 bytes)    │
└────────────────────────────┘
Effect: Bytes 1006-1010 sent
        Next available: 1011


CONNECTION CLOSE:
════════════════

FIN also consumes 1 sequence number!

Client initiates close:
┌────────────────────────────┐
│ FIN Flag: 1, ACK Flag: 1   │
│ Seq: 1011                  │ ← Current sequence
│ Ack: 5006                  │
└────────────────────────────┘
Effect: Seq 1011 consumed by FIN
        Next would be: 1012
```

### Sequence Number Increment Rules

```
┌─────────────────────────┬──────────────────┬─────────────────┐
│ Packet Type             │ Seq Increases?   │ By How Much?    │
├─────────────────────────┼──────────────────┼─────────────────┤
│ SYN                     │ Yes              │ +1              │
│ FIN                     │ Yes              │ +1              │
│ Data (N bytes)          │ Yes              │ +N              │
│ Pure ACK (no data)      │ No               │ 0               │
│ RST (Reset)             │ No               │ 0               │
└─────────────────────────┴──────────────────┴─────────────────┘
```

### Why Sequence Numbers Increase

1. **Byte-Level Tracking**: Each byte needs a unique identifier
   ```
   Sending "ABC":
   - 'A' is byte 1001
   - 'B' is byte 1002
   - 'C' is byte 1003
   ```

2. **Reordering**: If packets arrive out of order, receiver can sort them
   ```
   Packet 2 (Seq=1101) arrives before Packet 1 (Seq=1001)
   Receiver: "Wait, I'm missing 1001-1100, let me buffer 1101++"
   ```

3. **Loss Detection**: If ACK doesn't advance, sender knows packet was lost
   ```
   Send: Seq=1001
   Receive: Ack=1001 (duplicate ACK - something's wrong!)
   ```

4. **Duplicate Detection**: Same sequence number = duplicate packet
   ```
   Receive: Seq=1001 (first time)
   Receive: Seq=1001 (second time - duplicate, discard!)
   ```

### Sequence Number Wrapping

Since sequence numbers are 32-bit, they wrap around after 4GB:

```
Sequence at: 4,294,967,295 (max 32-bit value)
Next byte:   0 (wraps around)

TCP handles this correctly!

Example:
4,294,967,290 → 4,294,967,291 → ... → 4,294,967,295 → 0 → 1 → 2
```

### Practical Scenarios

#### Scenario 1: Perfect Delivery
```
Client                              Server
  |                                   |
  | Seq=1001, 100 bytes              |
  |--------------------------------->|
  |                                   | ✓ Received bytes 1001-1100
  |                                   |
  |              Ack=1101             |
  |<---------------------------------|
  |                                   |
  | Seq=1101, 50 bytes               |
  |--------------------------------->|
  |                                   | ✓ Received bytes 1101-1150
  |                                   |
  |              Ack=1151             |
  |<---------------------------------|
```

#### Scenario 2: Packet Loss and Retransmission
```
Client                              Server
  |                                   |
  | Seq=1001, 100 bytes              |
  |--------X (lost!)                 |
  |                                   |
  | (timeout, no ACK)                |
  |                                   |
  | Seq=1001, 100 bytes (RETRY)      |
  |--------------------------------->|
  |                                   | ✓ Received bytes 1001-1100
  |                                   |
  |              Ack=1101             |
  |<---------------------------------|
```

#### Scenario 3: Out-of-Order Delivery
```
Client                              Server
  |                                   |
  | Seq=1001, 100 bytes              |
  |---------\                         |
  |          \                        |
  | Seq=1101, 100 bytes              |
  |----------\-----------------------|
  |           \                       | ✓ Received 1101-1200 (buffered)
  |            \                      | ⏳ Waiting for 1001-1100
  |             \-------------------->|
  |                                   | ✓ Received 1001-1100
  |                                   | ✓ Combined with buffered data
  |              Ack=1201             |
  |<---------------------------------|
  |       (acknowledges all data)     |
```

#### Scenario 4: Simultaneous Data Transfer (Full Duplex)
```
Time  Client (Seq → / ← Ack)         Server (Seq → / ← Ack)
────  ──────────────────────         ──────────────────────
T1    Seq=1001, 100 bytes     →
T2                             →      Received, Ack=1101
T3                             ←      Seq=5001, 200 bytes
T4    Received, Ack=5201       ←
T5    Seq=1101, 50 bytes       →
      Ack=5201                 →      
T6                             →      Received both
                               ←      Ack=1151, Seq=5201, 150 bytes
T7    Received, Ack=5351       →

Final State:
Client sent: 150 bytes (1001-1150)
Server sent: 350 bytes (5001-5350)
Both sides acknowledged all data
```

### Common Mistakes and Misconceptions

❌ **Wrong**: "Sequence number is the packet number"
✅ **Right**: "Sequence number is the byte number"

❌ **Wrong**: "Each packet increments sequence by 1"
✅ **Right**: "Sequence increments by the number of bytes sent"

❌ **Wrong**: "ACKs increase sequence numbers"
✅ **Right**: "Pure ACKs don't change sequence numbers"

❌ **Wrong**: "Both sides use the same sequence number"
✅ **Right**: "Each direction has independent sequence numbers"

### Debugging with Sequence Numbers

When analyzing TCP with tools like Wireshark:

```bash
# Look for these patterns:

1. Retransmissions:
   Same Seq appearing multiple times

2. Out-of-order:
   Seq numbers not increasing monotonically

3. Lost packets:
   Gap in sequence numbers, then duplicate ACKs

4. Window full:
   Seq numbers stop advancing (flow control)
```

### Security Considerations

**TCP Sequence Number Prediction Attack:**
```
Old vulnerability (now fixed):
1. Attacker guesses next sequence number
2. Injects fake packet with correct Seq
3. Server accepts fake data

Modern fix:
- Random Initial Sequence Numbers
- Cryptographic randomness
- Makes prediction nearly impossible
```

### Summary Table

```
┌──────────────────────────────────────────────────────────────┐
│ Sequence Number Quick Reference                              │
├──────────────────────────────────────────────────────────────┤
│ Purpose:          Track bytes transmitted                    │
│ Size:             32 bits (0 to 4,294,967,295)              │
│ Initial Value:    Random (ISN)                               │
│ Increases When:   Sending data, SYN, or FIN                  │
│ Stays Same:       Pure ACKs, RST                             │
│ Per Connection:   Two independent sequences (client/server)  │
│ Used For:         Ordering, reliability, duplicate detection │
└──────────────────────────────────────────────────────────────┘
```

---

## ⚔️ TCP vs UDP: The Battle

### TCP: The Reliable Friend

**Analogy**: Sending a certified mail package
- You get tracking
- Delivery confirmation
- Package arrives intact and in order
- But it takes longer

**Code Example (Conceptual)**:
```
TCP sends: "Hello World"

Behind the scenes:
1. Establish connection (handshake)
2. Send "Hello" → Wait for ACK
3. Send " " → Wait for ACK
4. Send "World" → Wait for ACK
5. Close connection

Result: Receiver gets exactly "Hello World" in order
```

### UDP: The Fast Messenger

**Analogy**: Shouting across the room
- No confirmation needed
- Just send and hope they hear
- Fast but unreliable

**Code Example (Conceptual)**:
```
UDP sends: "Hello World"

Behind the scenes:
1. Send "Hello World" packet
2. Done!

Result: Receiver might get:
- "Hello World" ✓
- "Hell orld" (partial loss)
- Nothing (complete loss)
- "World Hello" (out of order - though rare in UDP)
```

### Real-World Decision Matrix

| Application | Protocol | Reason |
|------------|----------|--------|
| Web Browsing | TCP | Need complete, correct pages |
| Video Call | UDP | Speed > perfect quality |
| Email | TCP | Can't lose messages |
| Online Gaming | UDP | Real-time position updates |
| File Download | TCP | Need complete file |
| Live Sports Stream | UDP | Old frames aren't useful |
| Database Query | TCP | Need exact data |
| DNS Lookup | UDP | Quick, can retry if needed |

---

## 🔌 Socket Programming Concepts

### What is a Socket?

A **socket** is an endpoint for network communication. Think of it as a phone:
- Phone number = IP address + Port
- Making a call = Connecting
- Speaking/Listening = Sending/Receiving data

### Socket Types

1. **Stream Socket (TCP)**:
   - Connection-oriented
   - Reliable, ordered delivery
   - Like a phone call

2. **Datagram Socket (UDP)**:
   - Connectionless
   - Unreliable, unordered
   - Like sending postcards

### Socket Programming Flow

#### Server (TCP):
```
1. socket()      - Create a socket
2. bind()        - Bind to address:port
3. listen()      - Mark socket as passive (ready to accept)
4. accept()      - Wait for incoming connection
5. read/write()  - Communicate with client
6. close()       - Close connection
```

#### Client (TCP):
```
1. socket()      - Create a socket
2. connect()     - Connect to server
3. write/read()  - Communicate with server
4. close()       - Close connection
```

### Concurrent Server Patterns

#### 1. Thread-per-Connection (Go: Goroutine-per-Connection)
```
Server (Main Thread)
  ├── Goroutine 1 → Client 1
  ├── Goroutine 2 → Client 2
  └── Goroutine 3 → Client 3
```

**Pros**: Simple, isolated client handling
**Cons**: Resource usage with many connections

#### 2. Event-Driven (Select/Poll/Epoll)
```
Server (Single Thread)
  └── Event Loop
      ├── Check Client 1 (data available?)
      ├── Check Client 2 (data available?)
      └── Check Client 3 (data available?)
```

**Pros**: Handles many connections efficiently
**Cons**: More complex code

---

## 📨 Message Protocols

### Why Structured Messages?

Raw TCP provides a byte stream, not discrete messages. You need a protocol to frame messages.

### Protocol Options

#### 1. Length-Prefixed
```
[4 bytes: length][N bytes: data]

Example:
[0x00 0x00 0x00 0x0B]["Hello World"]
```

#### 2. Delimiter-Based
```
[data][delimiter]

Example:
"Hello World\n"
```

#### 3. Fixed-Length
```
[100 bytes: always]

Example:
"Hello World" + padding to 100 bytes
```

#### 4. Self-Describing (JSON, XML, Protocol Buffers)
```json
{
  "type": "message",
  "data": "Hello World"
}
```

**Our Implementation** uses JSON (self-describing) with newline delimiters:
```json
{"command":"ECHO","data":"Hello"}\n
```

---

## 🌍 Real-World Socket Applications

### 1. Web Servers
- **Technology**: HTTP over TCP
- **Sockets**: Listen on port 80 (HTTP) or 443 (HTTPS)
- **Example**: Nginx, Apache

### 2. Chat Applications
- **Technology**: TCP for reliability, sometimes WebSocket
- **Sockets**: Custom protocol over TCP
- **Example**: Slack, Discord

### 3. Multiplayer Games
- **Technology**: Often UDP for game state, TCP for critical events
- **Sockets**: UDP for position updates, TCP for chat/inventory
- **Example**: Fortnite, Call of Duty

### 4. Database Connections
- **Technology**: TCP
- **Sockets**: Custom protocol (MySQL protocol, PostgreSQL protocol)
- **Example**: PostgreSQL on port 5432

### 5. Email
- **Technology**: TCP
- **Sockets**: SMTP (port 25), IMAP (port 143), POP3 (port 110)
- **Example**: Mail servers

### 6. Video Streaming
- **Technology**: UDP (live) or TCP (on-demand)
- **Sockets**: RTMP, WebRTC (UDP), HLS (TCP)
- **Example**: YouTube, Twitch

### 7. IoT Devices
- **Technology**: MQTT over TCP, or CoAP over UDP
- **Sockets**: Lightweight protocols for constrained devices
- **Example**: Smart home devices

---

## 🎯 Best Practices

### 1. Error Handling
```go
conn, err := net.Dial("tcp", "localhost:8080")
if err != nil {
    log.Fatal(err) // Always handle errors!
}
defer conn.Close() // Always clean up
```

### 2. Timeouts
```go
conn.SetReadDeadline(time.Now().Add(30 * time.Second))
// Prevents hanging forever
```

### 3. Buffer Management
```go
buffer := make([]byte, 4096) // Reasonable buffer size
n, err := conn.Read(buffer)
data := buffer[:n] // Use only what was read
```

### 4. Graceful Shutdown
```go
// Handle signals for clean shutdown
sigChan := make(chan os.Signal, 1)
signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
<-sigChan
server.Shutdown()
```

### 5. Resource Limits
```go
// Limit concurrent connections
semaphore := make(chan struct{}, 100) // Max 100 connections
```

---

## 🐛 Common Pitfalls

### 1. Not Handling Partial Reads/Writes
```go
// ❌ Wrong
n, _ := conn.Read(buffer)

// ✅ Correct
for {
    n, err := conn.Read(buffer)
    if err != nil {
        break
    }
    process(buffer[:n])
}
```

### 2. Forgetting to Close Connections
```go
// ❌ Wrong
conn, _ := net.Dial("tcp", "localhost:8080")
// ... use conn ...
// Forgot to close!

// ✅ Correct
conn, _ := net.Dial("tcp", "localhost:8080")
defer conn.Close()
```

### 3. Blocking on Single Client
```go
// ❌ Wrong (Server handles one client at a time)
for {
    conn, _ := listener.Accept()
    handleClient(conn) // Blocks!
}

// ✅ Correct (Concurrent handling)
for {
    conn, _ := listener.Accept()
    go handleClient(conn) // Non-blocking
}
```

### 4. Not Validating Input
```go
// ❌ Wrong
var msg Message
json.Unmarshal(data, &msg)
process(msg) // What if unmarshal failed?

// ✅ Correct
var msg Message
if err := json.Unmarshal(data, &msg); err != nil {
    return err
}
if msg.Command == "" {
    return errors.New("invalid message")
}
process(msg)
```

---

## 🚀 Performance Considerations

### 1. Connection Pooling
Reuse connections instead of creating new ones:
```
Client → Connection Pool → Server
         ├── Conn 1 (idle)
         ├── Conn 2 (in use)
         └── Conn 3 (idle)
```

### 2. Buffering
Use bufio for efficient I/O:
```go
reader := bufio.NewReader(conn)
writer := bufio.NewWriter(conn)
```

### 3. Keep-Alive
Keep connections open for multiple requests:
```go
conn.SetKeepAlive(true)
conn.SetKeepAlivePeriod(30 * time.Second)
```

---

## 📚 Further Study Topics

1. **TLS/SSL**: Secure sockets
2. **WebSockets**: Full-duplex communication over HTTP
3. **HTTP/2 & HTTP/3**: Modern web protocols
4. **QUIC**: UDP-based transport protocol
5. **Load Balancing**: Distributing connections across servers
6. **NAT Traversal**: Dealing with firewalls and NAT
7. **Protocol Buffers**: Efficient serialization
8. **gRPC**: Modern RPC framework

---

Happy Learning! 🎓


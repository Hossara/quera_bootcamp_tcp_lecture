# TCP Sequence Numbers - Quick Reference Guide

## What is a Sequence Number?

A **TCP sequence number** is a 32-bit counter in the TCP header that uniquely identifies each byte of data transmitted over a TCP connection.

## Key Points

### 1. Purpose
- **Ordering**: Ensures data arrives in correct order
- **Reliability**: Detects missing or duplicated packets
- **Flow Control**: Tracks what's been sent and acknowledged
- **Byte Identification**: Each byte has a unique number

### 2. How It Works

```
Initial Sequence Number (ISN):
- Randomly chosen when connection starts
- Client picks one ISN (e.g., 1000)
- Server picks another ISN (e.g., 5000)
- Each direction has independent sequence numbers

Example:
Client ISN: 1000 → After handshake starts at 1001
Server ISN: 5000 → After handshake starts at 5001
```

### 3. When Does It Increase?

| Event          | Sequence Increase | Example                              |
|----------------|-------------------|--------------------------------------|
| **SYN flag**   | +1                | Seq=1000 → Next=1001                 |
| **FIN flag**   | +1                | Seq=1500 → Next=1501                 |
| **Data bytes** | +N bytes          | Seq=1001, Send 100 bytes → Next=1101 |
| **Pure ACK**   | 0 (no change)     | Seq=1101 stays 1101                  |

### 4. Why Does It Increase?

**Each byte needs a unique identifier:**
```
Sending "HELLO WORLD" (11 bytes):

Byte positions:
H = 1001
E = 1002
L = 1003
L = 1004
O = 1005
  = 1006  (space)
W = 1007
O = 1008
R = 1009
L = 1010
D = 1011

Next sequence: 1012
```

**Benefits:**
1. **Reordering**: If packet 2 arrives before packet 1, receiver can sort them
2. **Loss Detection**: Missing sequence numbers indicate lost packets
3. **Duplicate Detection**: Same sequence number = duplicate (discard)
4. **Acknowledgment**: Receiver says "I got everything up to byte X"

## Complete Example

### Three-Way Handshake

```
Client → Server: SYN
  Seq=1000
  (SYN consumes 1 number, next will be 1001)

Server → Client: SYN-ACK
  Seq=5000, Ack=1001
  (Server's SYN consumes 1, next will be 5001)
  (Ack=1001 means "I got your SYN at 1000, send 1001 next")

Client → Server: ACK
  Seq=1001, Ack=5001
  (Connection established)
```

### Data Transfer

```
Client sends 100 bytes:
  Seq=1001, Data=[100 bytes]
  Next sequence: 1001 + 100 = 1101

Server acknowledges:
  Ack=1101
  Meaning: "I received bytes 1001-1100, send 1101 next"

Client sends 50 more bytes:
  Seq=1101, Data=[50 bytes]
  Next sequence: 1101 + 50 = 1151

Server acknowledges:
  Ack=1151
  Meaning: "I received bytes 1101-1150, send 1151 next"
```

## Visualization

```
Time Flow:
═════════

T0: Client Seq=1000 (ISN chosen)
T1: After SYN: Seq=1001 (SYN consumed 1)
T2: Send 100 bytes → Seq=1001, Next=1101
T3: Send 50 bytes → Seq=1101, Next=1151
T4: Send 25 bytes → Seq=1151, Next=1176
T5: Send FIN → Seq=1176, Next=1177 (FIN consumed 1)

Total progression: 1000 → 1001 → 1101 → 1151 → 1176 → 1177
                   ISN   +1(SYN) +100   +50    +25    +1(FIN)
```

## Important Facts

1. **32-bit number**: Range is 0 to 4,294,967,295
2. **Wraps around**: After max value, returns to 0
3. **Bidirectional**: Client and server have separate sequences
4. **Random start**: ISN is random for security
5. **Byte-level**: Tracks individual bytes, not packets
6. **Cumulative ACK**: ACK number = "next byte I expect"

## Common Scenarios

### Packet Loss
```
Send: Seq=1001, 100 bytes
(Packet lost!)
Timeout → Retransmit: Seq=1001, 100 bytes (same sequence!)
Server ACKs: Ack=1101
```

### Out of Order
```
Send: Packet A (Seq=1001, 100 bytes) - delayed
Send: Packet B (Seq=1101, 100 bytes) - arrives first

Server receives B: Buffers it, sends Ack=1001 (waiting for A)
Server receives A: Combines both, sends Ack=1201
```

### Pure ACK
```
Server sends: Seq=5001, Ack=1101 (no data)
Server's sequence: Still 5001 (unchanged)

Server sends: Seq=5001, Ack=1151 (no data)
Server's sequence: Still 5001 (unchanged)

Only data/SYN/FIN increase sequence!
```

## Quick Mental Model

Think of sequence numbers as **page numbers in a book**:

- Each page (byte) has a unique number
- Pages can arrive out of order - you reorder them
- If page 5 is missing, you know exactly which one to request
- You can tell if you got the same page twice
- You acknowledge by saying "I have up to page 100, send page 101 next"

## For Debugging

When analyzing with Wireshark or tcpdump:

```
Look for:
- Same Seq appearing twice = Retransmission
- Seq gaps = Packet loss
- Seq not increasing = Out of order or problem
- Duplicate ACKs = Receiver waiting for missing packet
```

## Summary

**What**: 32-bit byte counter in TCP header

**When it increases**:
- SYN flag: +1
- FIN flag: +1  
- Data: +number of bytes
- ACK only: +0

**Why it increases**:
- Track every byte uniquely
- Enable ordering, reliability, duplicate detection

**How it works**:
- Random start (ISN)
- Increments by bytes sent
- Each direction independent
- Wraps around at 2^32

---

For complete details and examples, see **CONCEPTS.md** in this repository.


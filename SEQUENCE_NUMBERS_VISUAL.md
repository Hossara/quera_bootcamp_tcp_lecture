# TCP Sequence Number Visual Guide

This document provides visual representations of how TCP sequence numbers work.

## Visual 1: Sequence Number Timeline

```
Connection Lifecycle with Sequence Numbers:
═══════════════════════════════════════════════════════════════

                    CLIENT SIDE                   |            SERVER SIDE
                                                  |
ISN Selection       Seq = 1000 (random)           |           Seq = 5000 (random)
                         │                        |                │
                         │                        |                │
HANDSHAKE           ┌────▼────┐                   |           ┌────▼────┐
                    │  1000   │ SYN              →│→          │  5000   │
                    └────┬────┘                   |           └────┬────┘
                         │                        |                │
                    ┌────▼────┐                   |           ┌────▼────┐
                    │  1001   │                  ←│←  SYN-ACK │  5001   │
                    └────┬────┘                   |           └────┬────┘
                         │                        |                │
                    ┌────▼────┐  ACK             →│→               │
                    │  1001   │ (no data)         |                │
                    └────┬────┘                   |                │
                         │                        |                │
DATA TRANSFER            │                        |                │
                    ┌────▼────┐                   |                │
                    │  1001   │ 100 bytes        →│→               │
                    └────┬────┘                   |                │
                         │                        |                │
                    ┌────▼────┐                  ←│←  ACK=1101     │
                    │  1101   │                   |                │
                    └────┬────┘                   |                │
                         │                        |           ┌────▼────┐
                         │                       ←│←  200 bytes│  5001   │
                         │                        |           └────┬────┘
                    ┌────▼────┐                   |                │
                    │  1101   │ ACK=5201         →│→          ┌────▼────┐
                    └────┬────┘                   |           │  5201   │
                         │                        |           └─────────┘
                    ┌────▼────┐                   |
                    │  1101   │ 50 bytes         →│→
                    └────┬────┘                   |
                         │                        |
                    ┌────▼────┐                  ←│←  ACK=1151
                    │  1151   │                   |
                    └────┬────┘                   |
                         │                        |
CONNECTION CLOSE    ┌────▼────┐                   |
                    │  1151   │ FIN              →│→
                    └────┬────┘                   |
                         │                        |
                    ┌────▼────┐                  ←│←  ACK=1152
                    │  1152   │                   |
                    └─────────┘                   |
```

## Visual 2: Byte-Level View

```
Sending "HELLO WORLD" (11 bytes):

Memory Buffer:
┌───┬───┬───┬───┬───┬───┬───┬───┬───┬───┬───┐
│ H │ E │ L │ L │ O │   │ W │ O │ R │ L │ D │
└───┴───┴───┴───┴───┴───┴───┴───┴───┴───┴───┘

Sequence Numbers (Starting at 1001):
┌─────┬─────┬─────┬─────┬─────┬─────┬─────┬─────┬─────┬─────┬─────┐
│1001 │1002 │1003 │1004 │1005 │1006 │1007 │1008 │1009 │1010 │1011 │
└─────┴─────┴─────┴─────┴─────┴─────┴─────┴─────┴─────┴─────┴─────┘

Packet 1: Seq=1001, Data="HELLO" (5 bytes)
┌─────────────────────┐
│ Seq: 1001           │ 
│ Data: H E L L O     │
│ Bytes: 1001-1005    │
└─────────────────────┘
Next Seq: 1006

Packet 2: Seq=1006, Data=" WORLD" (6 bytes)
┌─────────────────────┐
│ Seq: 1006           │
│ Data:   W O R L D   │
│ Bytes: 1006-1011    │
└─────────────────────┘
Next Seq: 1012
```

## Visual 3: Lost Packet Scenario

```
Time flows downward ↓

Client                              Server
Seq=1001                            Seq=5001
  │                                   │
  │                                   │
  ├─[Packet A: Seq=1001, 100 bytes]─→ │
  │                                   ├─ Received ✓
  │                                   │  Buffer: 1001-1100
  │                                   │
  │←────────[ACK=1101]────────────────┤
  │                                   │
  │                                   │
  ├─[Packet B: Seq=1101, 100 bytes]─X │  ← LOST!
  │                                   │
  │                                   │
  │  (Timeout waiting for ACK)        │
  │                                   │
  ├─[Packet B: Seq=1101, 100 bytes]─→ │  ← RETRANSMIT
  │  (Same sequence number!)          │
  │                                   ├─ Received ✓
  │                                   │  Buffer: 1001-1200
  │                                   │
  │←────────[ACK=1201]────────────────┤
  │                                   │
  │                                   │
Seq=1201                            Seq=5001
```

## Visual 4: Out-of-Order Delivery

```
Time flows downward ↓

Client                              Server
  │                                   │
  │                                   │
  ├─[Packet A: Seq=1001, 100 bytes]─┐││
  │                                  ││  ← Delayed in network
  │                                  ││
  ├─[Packet B: Seq=1101, 100 bytes]┐│││
  │                                 │││
  │                                 │││
  │                                 │└┼─→ Arrives SECOND
  │                                 │ │   Buffer: 1001-1100 ✓
  │                                 │ │
  │                                 └─┼─→ Arrives FIRST
  │                                   │   Out-of-order!
  │                                   │   Buffer: 1101-1200
  │                                   │   Waiting for: 1001-1100
  │                                   │
  │←────────[ACK=1001]────────────────┤
  │  "Still need 1001-1100"           │
  │                                   │
  │  (Packet A arrives)               │
  │                                   ├─ Now have both!
  │                                   │  Combine: 1001-1200 ✓
  │                                   │
  │←────────[ACK=1201]────────────────┤
  │  "Got everything!"                │
```

## Visual 5: Bidirectional Data Flow

```
Time    Client                          Server
═════   ══════                          ══════
        Seq=1001                        Seq=5001
        Ack=5001                        Ack=1001
          │                               │
T1        ├──[Seq=1001, 100B]────────────→│
          │                               │ Received bytes 1001-1100
T2        │                               │
          │←──────────[Ack=1101]──────────┤
          │                               │
T3        │←──[Seq=5001, 200B]────────────┤
          │ Received bytes 5001-5200      │
T4        │                               │
          ├────────[Ack=5201]────────────→│
          │                               │
T5        ├──[Seq=1101, 50B, Ack=5201]───→│
          │                               │ Received bytes 1101-1150
T6        │                               │
          │←─[Ack=1151, Seq=5201, 150B]───┤
          │ Received bytes 5201-5350      │
T7        │                               │
          ├────────[Ack=5351]────────────→│
          │                               │

Final:    Seq=1151                      Seq=5351
          Ack=5351                      Ack=1151

Summary:
- Client sent: 150 bytes (1001→1151)
- Server sent: 350 bytes (5001→5351)
- Both directions acknowledged
```

## Visual 6: Sequence Number Wrap-Around

```
32-bit maximum: 4,294,967,295

Near the end:
┌──────────────┬──────────────┬──────────────┬──────────────┬──────────────┐
│4,294,967,291 │4,294,967,292 │4,294,967,293 │4,294,967,294 │4,294,967,295 │
└──────────────┴──────────────┴──────────────┴──────────────┴──────────────┘
                                                                     │
                                                                     │ Wrap!
                                                                     ▼
┌──────────────┬──────────────┬──────────────┬──────────────┬──────────────┐
│      0       │      1       │      2       │      3       │      4       │
└──────────────┴──────────────┴──────────────┴──────────────┴──────────────┘

Example:
Send 10 bytes starting at Seq=4,294,967,293
Bytes occupy: 4,294,967,293, 294, 295, 0, 1, 2, 3, 4, 5, 6
Next Seq: 7
ACK will be: 7
```

## Visual 7: TCP Segment Structure

```
┌────────────────────────────────────────────────────────────┐
│                     TCP SEGMENT                            │
├────────────────────────────────────────────────────────────┤
│                                                            │
│  ┌──────────────────────────────────────────────────────┐  │
│  │                 TCP HEADER (20+ bytes)               │  │
│  ├──────────────────────────────────────────────────────┤  │
│  │ Source Port (16 bits)    │ Dest Port (16 bits)       │  │
│  ├──────────────────────────────────────────────────────┤  │
│  │          SEQUENCE NUMBER (32 bits)                   │  │ ← This!
│  ├──────────────────────────────────────────────────────┤  │
│  │         ACKNOWLEDGMENT NUMBER (32 bits)              │  │ ← And this!
│  ├──────────────────────────────────────────────────────┤  │
│  │ Offset  │ Flags (SYN,ACK,FIN) │ Window Size          │  │
│  ├──────────────────────────────────────────────────────┤  │
│  │ Checksum                    │ Urgent Pointer         │  │
│  └──────────────────────────────────────────────────────┘  │
│                                                            │
│  ┌──────────────────────────────────────────────────────┐  │
│  │              DATA / PAYLOAD                          │  │
│  │              (Variable length)                       │  │
│  └──────────────────────────────────────────────────────┘  │
│                                                            │
└────────────────────────────────────────────────────────────┘

Example values:
┌────────────────────────────────────────┐
│ Sequence Number:     1,458,234,567     │ ← Identifies byte position
│ Acknowledgment:      2,891,456,789     │ ← Confirms received bytes
│ Data Length:         1,460 bytes       │
│ Next Sequence:       1,458,236,027     │ ← 567 + 1460 = 2027
└────────────────────────────────────────┘
```

## Visual 8: State Diagram

```
TCP Connection States with Sequence Numbers:

        START
          │
          ├─── Pick Random ISN ───┐
          │                       │
     [CLOSED]                     │
          │                       │
          │ Active Open           │
          ▼                       │
    [SYN-SENT]                    │
    Seq=ISN                       │
          │                       │
          │ Receive SYN-ACK       │
          │ Ack=ISN+1             │
          ▼                       │
  [ESTABLISHED]                   │
  Seq=ISN+1                       │
  Can send/receive data           │
  Seq increases with data         │
          │                       │
          │ Close                 │
          ▼                       │
   [FIN-WAIT]                     │
   Seq=Current+1 (FIN)            │
          │                       │
          │ Receive ACK           │
          ▼                       │
     [CLOSED]                     │
     Connection ends              │
                                  │
                            Passive Open
                                  │
                                  ▼
                             [LISTEN]
                                  │
                         Receive SYN
                                  │
                                  ▼
                            [SYN-RECEIVED]
                            Send SYN-ACK
                            Seq=Server_ISN
                                  │
                         Receive ACK
                                  │
                                  └──────→ [ESTABLISHED]
```

## Key Formulas

```
┌─────────────────────────────────────────────────────────────┐
│ Sequence Number Calculations                                │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│ 1. Next Sequence Number:                                    │
│    Next_Seq = Current_Seq + Data_Length                     │
│                                                             │
│    Example: Seq=1001, Data=100 bytes                        │
│             Next_Seq = 1001 + 100 = 1101                    │
│                                                             │
│ 2. Acknowledgment Number:                                   │
│    ACK = Last_Received_Seq + 1                              │
│                                                             │
│    Example: Received up to byte 1100                        │
│             ACK = 1100 + 1 = 1101                           │
│             (Next expected byte is 1101)                    │
│                                                             │
│ 3. Bytes Acknowledged:                                      │
│    Bytes_Acked = ACK_Number - Initial_Seq                   │
│                                                             │
│    Example: Initial_Seq=1001, ACK=1151                      │
│             Bytes_Acked = 1151 - 1001 = 150 bytes           │
│                                                             │
│ 4. Outstanding Bytes (Unacknowledged):                      │
│    Outstanding = Current_Seq - Last_ACK_Received            │
│                                                             │
│    Example: Sent up to 1201, received ACK=1101              │
│             Outstanding = 1201 - 1101 = 100 bytes           │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

---

These visualizations complement the detailed explanations in CONCEPTS.md and help understand how TCP sequence numbers work in practice.


# Learning Redis - Step by Step Implementation

## What is Redis?

Redis (Remote Dictionary Server) is an in-memory data structure store that can be used as a database, cache, and message broker. It supports various data structures like strings, hashes, lists, sets, and more.

## Our Learning Journey

### Phase 1: Basic Server & Protocol

- [ ] TCP Server setup
- [ ] RESP (Redis Serialization Protocol) parser
- [ ] Basic command handling

### Phase 2: Core Data Structures

- [ ] String operations (GET, SET, DEL)
- [ ] Hash operations (HGET, HSET, HDEL)
- [ ] List operations (LPUSH, RPUSH, LPOP, RPOP)
- [ ] Set operations (SADD, SREM, SMEMBERS)

### Phase 3: Advanced Features

- [ ] Expiration (TTL, EXPIRE)
- [ ] Persistence (RDB snapshots)
- [ ] Pub/Sub messaging
- [ ] Transactions (MULTI, EXEC)

### Phase 4: Performance & Scalability

- [ ] Memory management
- [ ] Connection pooling
- [ ] Basic clustering concepts

## Architecture Overview

```
┌─────────────────┐     ┌─────────────────┐     ┌─────────────────┐
│   Redis Client  │────▶│   TCP Server    │────▶│  Command Router │
└─────────────────┘     └─────────────────┘     └─────────────────┘
                                                          │
                                                          ▼
┌─────────────────┐     ┌─────────────────┐     ┌─────────────────┐
│   Data Store    │◀────│   Storage       │◀────│  Command Handler│
│   (In-Memory)   │     │   Engine        │     └─────────────────┘
└─────────────────┘     └─────────────────┘
```

## Technology Stack

- **Language**: Go (for performance and excellent concurrency)
- **Networking**: TCP sockets with goroutines
- **Data Storage**: In-memory maps with sync.RWMutex
- **Protocol**: RESP (Redis Serialization Protocol)
- **Concurrency**: Goroutines and channels

## Why Go for Learning Redis?

1. **Performance**: Compiled language with near-C performance
2. **Concurrency**: Built-in goroutines make handling multiple clients easy
3. **Memory Safety**: Garbage collected, focus on logic not memory management
4. **Network Programming**: Excellent standard library for TCP servers
5. **Readability**: Clean syntax helps understand Redis concepts clearly

Let's start building!
# redis-clone

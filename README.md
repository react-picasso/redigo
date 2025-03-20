# 📌 RediGo

This is a Redis-like in-memory key-value store implemented in Go. It supports basic Redis commands, including `PING`, `ECHO`, `SET`, `GET`, `CONFIG GET`, `KEYS *`, and `SAVE`. It also features RDB persistence, allowing data to be saved and loaded from disk.

## 🚀 Features Implemented

### ✔️ Basic Commands
* ✅ `PING` → Responds with `PONG`
* ✅ `ECHO <message>` → Returns the same message
* ✅ `SET <key> <value>` → Stores a key-value pair
* ✅ `GET <key>` → Retrieves a value by key
* ✅ `KEYS "*"` → Returns all keys in the database

### ✔️ Advanced Features
* ✅ Persistence with RDB (`SAVE`) → Saves key-value pairs to an .rdb file
* ✅ Loading from RDB → Reads keys from an RDB file on startup
* ✅ `CONFIG GET dir/dbfilename` → Fetches RDB file location details
* ✅ Thread-safe key storage → Uses `sync.RWMutex` for safe concurrent access
* ✅ Automatic key expiry (PX option) → Deletes keys after a specified time

## 📂 Project Structure

```
- redis-clone/
    - cmd/
        - server/
            - main.go # Entrypoint for the server
        - internal/
            - logger/
                - logger.go # Handles logging
            - resp/
                - parser.go # Parses Redis RESP protocol
            - server/
                - tcp.go # Manages TCP connections
                - handler.go # Processes Redis commands
                - store.go # In-memory KV store
                - config.go # Handles config parameters
                - rdb.go # Implements RDB persistence
    - go.mod
    - go.sum
```

## 🛠️ Installation & Setup

### 1. Install Go

Make sure you have Go 1.18+ installed.

```
go version
```

### 2. Clone the Repository

```
git clone https://github.com/react-picasso/redis-clone-go.git
```

### 3. Run the Server

```
go run cmd/server/main.go --dir /tmp/redis-files --dbfilename dump.rdb
```

* The server will start on port `6379`.
* If an RDB file exists, it will load saved data.

## 📌 Usage

All the commands can be tested using `nc`:

### 1. Start netcat

```
nc localhost 6379
```

### 2. Test any command normally like you would with `redis-cli`

```
PING # Response PONG
SET foo bar # Response OK
GET foo # Response bar
```

## 📂 How RDB works

1. On startup, the server loads data from an RDB file (if available).
2. On `SAVE`, the server writes all keys to an RDB file in binary format.
3. The RDB format follows Redis serialization (header, metadata, key-value pairs, checksum).

## 🚀 Next Steps

🔷 Improve RDB persistence

🔷 Implement other features like replication, streams, and transactions

🔷 Improve RESP parsed responses
package server

import (
	"fmt"
	"net"
	"strings"
)

var store = NewStore()

// HandleCommand processes a parsed Redis command
func HandleCommand(command []string, conn net.Conn) {
	if len(command) == 0 {
		return
	}

	cmd := strings.ToUpper(command[0])

	switch cmd {
	case "PING":
		conn.Write([]byte("+PONG\r\n"))
	case "ECHO":
		if len(command) < 2 {
			conn.Write([]byte("-ERR missing argument\r\n"))
		} else {
			response := fmt.Sprintf("$%d\r\n%s\r\n", len(command[1]), command[1])
			conn.Write([]byte(response))
		}
	case "SET":
		if len(command) < 3 {
			conn.Write([]byte("-ERR wrong number of arguments for 'SET'\r\n"))
		} else {
			store.Set(command[1], command[2])
			conn.Write([]byte("+OK\r\n"))
		}
	case "GET":
		if len(command) < 2 {
			conn.Write([]byte("-ERR wrong number of arguments for 'GET'\r\n"))
		} else {
			value, exists := store.Get(command[1])
			if !exists {
				conn.Write([]byte("$-1\r\n"))
			} else {
				response := fmt.Sprintf("$%d\r\n%s\r\n", len(value), value)
				conn.Write([]byte(response))
			}
		}
	default:
		conn.Write([]byte("-ERR unknown command\r\n"))
		conn.Write([]byte("+OK\r\n"))
	}
}

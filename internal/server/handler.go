package server

import (
	"fmt"
	"net"
	"strconv"
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
			key := command[1]
			value := command[2]
			px := 0

			if len(command) >= 5 && strings.ToUpper(command[3]) == "PX" {
				if pxVal, err := strconv.Atoi(command[4]); err == nil {
					px = pxVal
				}
			}
			store.Set(key, value, px)
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
	case "CONFIG":
		if len(command) < 3 || strings.ToUpper(command[1]) != "GET" {
			conn.Write([]byte("-ERR invalid CONFIG command\r\n"))
		} else {
			param := strings.ToLower(command[2])
			var key, value string

			switch param {
			case "dir":
				key = "dir"
				value = ServerConfig.Dir
			case "dbfilename":
				key = "dbfilename"
				value = ServerConfig.DBFilename
			default:
				conn.Write([]byte("$-1\r\n")) // Null bulk string for unknown config keys
				return
			}

			response := fmt.Sprintf("*2\r\n$%d\r\n%s\r\n$%d\r\n%s\r\n", len(key), key, len(value), value)
			conn.Write([]byte(response))
		}
	case "KEYS":
		if len(command) < 2 || command[1] != "*" {
			conn.Write([]byte("-ERR invalid KEYS pattern\r\n"))
		} else {
			keys := store.GetAllKeys()
			resp := fmt.Sprintf("*%d\r\n", len(keys))
			for _, key := range keys {
				resp += fmt.Sprintf("$%d\r\n%s\r\n", len(key), key)
			}
			conn.Write([]byte(resp))
		}
	case "SAVE":
		err := SaveRDB()
		if err != nil {
			conn.Write([]byte(fmt.Sprintf("-ERR %s\r\n", err)))
		} else {
			conn.Write([]byte("+OK\r\n"))
		}
	default:
		conn.Write([]byte("-ERR unknown command\r\n"))
	}
}

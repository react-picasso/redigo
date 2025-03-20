package server

import (
	"bufio"
	"net"
	"strings"

	"github.com/react-picasso/redigo/internal/logger"
)

const port = ":6379"

func handleConnection(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	for {
		// Read client input
		msg, err := reader.ReadString('\n')
		if err != nil {
			logger.Logger.Println("Client disconnected")
			return
		}

		// Trim whitespace and print received message
		cmds := strings.Split(strings.TrimSpace(msg), "\n")

		for _, cmd := range cmds {
			cmd = strings.TrimSpace(cmd)
			if cmd == "" {
				continue
			}

			logger.Logger.Println("Received command:", cmd)

			response := "+PONG\r\n"
			_, err = conn.Write([]byte(response))
			if err != nil {
				logger.Logger.Println("Error writing response:", err)
				return
			}
		}
	}
}

func StartServer() {
	lsnr, err := net.Listen("tcp", port)
	if err != nil {
		logger.Logger.Fatalf("Error starting server: %v", err)
	}
	defer lsnr.Close()

	logger.Logger.Println("Server started on port", port)

	for {
		conn, err := lsnr.Accept()
		if err != nil {
			logger.Logger.Println("Connection error:", err)
			continue
		}

		logger.Logger.Println("New client connected")
		go handleConnection(conn)
	}
}

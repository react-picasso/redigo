package server

import (
	"bufio"
	"net"

	"github.com/react-picasso/redigo/internal/logger"
	"github.com/react-picasso/redigo/internal/resp"
)

const port = ":6379"

func handleConnection(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	for {
		// Parse RESP command
		command, err := resp.ParseRESP(reader)
		if err != nil {
			logger.Logger.Println("Client disconnected or invalid data:", err)
			return
		}

		logger.Logger.Println("Received command:", command)
		HandleCommand(command, conn)
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

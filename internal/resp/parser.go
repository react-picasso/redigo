package resp

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
)

// ParseRESP parses the Redis protocol and extracts commands and arguments
func ParseRESP(reader *bufio.Reader) ([]string, error) {
	// Read first line
	line, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}
	line = strings.TrimSpace(line)

	// If the input is raw text (e.g., from netcat), split manually
	if !strings.HasPrefix(line, "*") {
		return strings.Fields(line), nil
	}

	// Handle proper RESP-encoded input
	numElements, err := strconv.Atoi(line[1:])
	if err != nil {
		return nil, fmt.Errorf("invalid RESP array length")
	}

	commands := make([]string, numElements)
	for i := 0; i < numElements; i++ {
		_, err := reader.ReadString('\n') // Skip the bulk string length line
		if err != nil {
			return nil, err
		}

		// Read the actual command/argument
		data, err := reader.ReadString('\n')
		if err != nil {
			return nil, err
		}
		commands[i] = strings.TrimSpace(data)
	}

	return commands, nil
}

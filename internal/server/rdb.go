package server

import (
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func LoadRDB() {
	rdbPath := filepath.Join(ServerConfig.Dir, ServerConfig.DBFilename)
	file, err := os.Open(rdbPath)
	if err != nil {
		fmt.Println("RDB file not found, starting with an empty database.")
		return
	}
	defer file.Close()

	// Header: Magic string + version
	header := make([]byte, 9)
	_, err = file.Read(header)
	if err != nil || !strings.HasPrefix(string(header), "REDIS0011") {
		fmt.Println("Invalid RDB file format")
		return
	}

	// Start reading sections
	for {
		// Read next byte
		var b [1]byte
		_, err := file.Read(b[:])
		if err != nil {
			break
		}

		switch b[0] {
		case 0xFE: // Database selector
			continue
		case 0xFB: // Hash table size info
			file.Read(b[:]) // Read size, ignore
			file.Read(b[:]) // Read expires, ignore
		case 0xFC, 0xFD: // Expiry timestamps (ignore for now)
			expiry := make([]byte, 8)
			file.Read(expiry)
		case 0x00: // Key-Value pair (string type)
			key, _ := readString(file)
			value, _ := readString(file)
			store.Set(key, value, 0) // No expiry for now
		case 0xFF: // End of file
			return
		}
	}
}

func readString(file *os.File) (string, error) {
	size, err := readSize(file)
	if err != nil {
		return "", err
	}

	data := make([]byte, size)
	_, err = file.Read(data)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

// readSize decodes a size-encoded value from an RDB file
func readSize(file *os.File) (int, error) {
	var b [1]byte
	_, err := file.Read(b[:])
	if err != nil {
		return 0, err
	}

	switch b[0] >> 6 {
	case 0b00: // 6-bit size
		return int(b[0] & 0x3F), nil
	case 0b01: // 14-bit size
		var nextByte [1]byte
		_, err := file.Read(nextByte[:])
		if err != nil {
			return 0, err
		}
		return int(b[0]&0x3F)<<8 | int(nextByte[0]), nil
	case 0b10: // 32-bit size
		var sizeBytes [4]byte
		_, err := file.Read(sizeBytes[:])
		if err != nil {
			return 0, err
		}
		return int(binary.BigEndian.Uint32(sizeBytes[:])), nil
	}
	return 0, fmt.Errorf("invalid size encoding")
}

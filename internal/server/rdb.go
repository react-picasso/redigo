package server

import (
	"encoding/binary"
	"fmt"
	"hash/crc64"
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

func SaveRDB() error {
	rdbPath := filepath.Join(ServerConfig.Dir, ServerConfig.DBFilename)
	file, err := os.Create(rdbPath)
	if err != nil {
		return fmt.Errorf("failed to create RDB file: %v", err)
	}
	defer file.Close()

	// 1. Write the header
	file.Write([]byte("REDIS0011"))

	// 2. Write metadata
	writeMetadata(file, "redis-ver", "6.0.16")

	// 3. Write the database section
	file.Write([]byte{0xFE, 0x00})                        // Database selector (DB 0)
	file.Write([]byte{0xFB, byte(len(store.data)), 0x00}) // Hash table sizes

	// 4. Write key-value pairs
	for key, value := range store.data {
		file.Write([]byte{0x00}) // String type
		writeString(file, key)
		writeString(file, value)
	}

	// 5. Write EOF marker and checksum
	file.Write([]byte{0xFF}) // End of file
	writeChecksum(file)

	return nil
}

// writeMetadata writes a metadata key-value pair
func writeMetadata(file *os.File, key, value string) {
	file.Write([]byte{0xFA}) // Metadata section
	writeString(file, key)
	writeString(file, value)
}

// writeString writes a size-encoded string
func writeString(file *os.File, str string) {
	size := len(str)
	if size <= 63 {
		file.Write([]byte{byte(size)})
	} else if size <= 16383 {
		file.Write([]byte{byte(0x40 | (size >> 8)), byte(size & 0xFF)})
	} else {
		file.Write([]byte{0x80, byte(size >> 24), byte(size >> 16), byte(size >> 8), byte(size)})
	}
	file.Write([]byte(str))
}

// writeChecksum writes a CRC64 checksum of the entire file
func writeChecksum(file *os.File) {
	file.Seek(0, 0)
	data, _ := os.ReadFile(file.Name())
	crc := crc64.Checksum(data, crc64.MakeTable(crc64.ISO))
	binary.Write(file, binary.LittleEndian, crc)
}

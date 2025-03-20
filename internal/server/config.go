package server

import "flag"

type Config struct {
	Dir        string
	DBFilename string
}

// Global instance of config
var ServerConfig = &Config{}

// ParseFlags parses command-line arguments for configuration
func ParseFlags() {
	dir := flag.String("dir", "/tmp/redis-data", "Directory to store RDB files")
	dbfilename := flag.String("dbfilename", "dump.rdb", "RDB filename")

	flag.Parse()

	ServerConfig.Dir = *dir
	ServerConfig.DBFilename = *dbfilename
}

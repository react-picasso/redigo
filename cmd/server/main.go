package main

import "github.com/react-picasso/redigo/internal/server"

func main() {
	server.ParseFlags()
	server.LoadRDB()
	server.StartServer()
}

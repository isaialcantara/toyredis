package main

import "github.com/isaialcantara/toyredis/internal/server"

func main() {
	server := server.New(":6379")
	err := server.Start()
	if err != nil {
		panic(err)
	}
}

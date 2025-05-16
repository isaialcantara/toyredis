package main

import (
	"github.com/isaialcantara/toyredis/internal/server"
)

func main() {
	myServer := server.New(":6379")
	err := myServer.Start()
	if err != nil {
		panic(err)
	}
}

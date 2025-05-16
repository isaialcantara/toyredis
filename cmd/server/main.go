package main

import (
	"github.com/isaialcantara/toyredis/internal/command"
	"github.com/isaialcantara/toyredis/internal/server"
)

func main() {
	dispatcher := command.NewBasicDispatcher()
	myServer := server.New(":6379", dispatcher)
	err := myServer.Start()
	if err != nil {
		panic(err)
	}
}

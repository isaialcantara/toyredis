package main

import (
	"github.com/isaialcantara/toyredis/internal/command"
	"github.com/isaialcantara/toyredis/internal/server"
	"github.com/isaialcantara/toyredis/internal/storage"
)

func main() {
	store := storage.NewMapStore()
	dispatcher := command.NewBasicDispatcher(store)
	myServer := server.New(":6379", dispatcher)

	err := myServer.Start()
	panic(err)
}

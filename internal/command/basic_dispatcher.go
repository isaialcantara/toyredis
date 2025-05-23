package command

import (
	"log"
	"strings"

	"github.com/isaialcantara/toyredis/internal/resp"
)

type command struct {
	handler func(resp.BulkArray) []byte
	name    string
	summary string
}

type BasicDispatcher struct {
	commands map[string]*command
}

var _ Dispatcher = (*BasicDispatcher)(nil)

func NewBasicDispatcher() *BasicDispatcher {
	return &BasicDispatcher{
		commands: map[string]*command{
			"PING": {name: "ping", summary: "Pings the server.", handler: pingHandler},
			"ECHO": {name: "echo", summary: "Echoes the givem message.", handler: echoHandler},
		},
	}
}

func (d *BasicDispatcher) Dispatch(bulkArray resp.BulkArray) []byte {
	if len(bulkArray) < 1 {
		return ErrCommandEmpty.ToRESP()
	}

	cmd, args := d.findCommand(bulkArray)
	if cmd == nil {
		log.Printf("invalid command: %q", bulkArray.ToRESP())
		return ErrCommandInvalid.ToRESP()
	}
	return cmd.handler(args)
}

func (d *BasicDispatcher) findCommand(bulkArray resp.BulkArray) (*command, resp.BulkArray) {
	var path string
	var cmd *command
	var argsStart int
	for i, bulk := range bulkArray {
		part := strings.ToUpper(string(bulk))
		if i == 0 {
			path = part
		} else {
			path += " " + part
		}

		if value, ok := d.commands[path]; ok {
			cmd = value
			argsStart = i + 1
			continue
		}
		break
	}

	return cmd, bulkArray[argsStart:]
}

func pingHandler(args resp.BulkArray) []byte {
	if len(args) == 0 {
		return resp.SimpleString("PONG").ToRESP()
	}
	return echoHandler(args)
}

func echoHandler(args resp.BulkArray) []byte {
	if len(args) != 1 {
		return ErrCommandArgsNumber.ToRESP()
	}
	return resp.BulkString(string(args[0])).ToRESP()
}

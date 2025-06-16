package command

import (
	"log"
	"strings"

	"github.com/isaialcantara/toyredis/internal/resp"
	"github.com/isaialcantara/toyredis/internal/storage"
)

type command struct {
	handler func(resp.BulkArray) []byte
	name    string
	summary string
}

type BasicDispatcher struct {
	commands map[string]*command
	store    storage.StringStorage
}

var _ Dispatcher = (*BasicDispatcher)(nil)

func NewBasicDispatcher(store storage.StringStorage) *BasicDispatcher {
	d := &BasicDispatcher{store: store}

	commands := map[string]*command{
		"PING": {name: "ping", summary: "Pings the server.", handler: d.pingHandler},
		"ECHO": {name: "echo", summary: "Returns message.", handler: d.echoHandler},
		"GET":  {name: "get", summary: "Get the value of key (string only).", handler: d.getHandler},
		"SET":  {name: "set", summary: "Set the value of key (string only).", handler: d.setHandler},
		"DEL":  {name: "del", summary: "Deletes the value at the given key.", handler: d.delHandler},
	}

	d.commands = commands
	return d
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

func (d *BasicDispatcher) pingHandler(args resp.BulkArray) []byte {
	if len(args) == 0 {
		return resp.SimpleString("PONG").ToRESP()
	}
	return d.echoHandler(args)
}

func (d *BasicDispatcher) echoHandler(args resp.BulkArray) []byte {
	if len(args) != 1 {
		return ErrCommandArgsNumber.ToRESP()
	}
	return resp.BulkString(string(args[0])).ToRESP()
}

func (d *BasicDispatcher) getHandler(args resp.BulkArray) []byte {
	if len(args) != 1 {
		return ErrCommandArgsNumber.ToRESP()
	}

	val, err := d.store.Get(string(args[0]))
	if err != nil {
		return ErrCommandFailed.ToRESP()
	}

	return resp.BulkString(val).ToRESP()
}

func (d *BasicDispatcher) setHandler(args resp.BulkArray) []byte {
	if len(args) != 2 {
		return ErrCommandArgsNumber.ToRESP()
	}

	if err := d.store.Set(string(args[0]), args[1]); err != nil {
		return ErrCommandFailed.ToRESP()
	}

	return resp.SimpleString("OK").ToRESP()
}

func (d *BasicDispatcher) delHandler(args resp.BulkArray) []byte {
	if len(args) < 1 {
		return ErrCommandArgsNumber.ToRESP()
	}

	deleted := 0
	for _, key := range args {
		wasDeleted, err := d.store.Del(string(key))
		if err != nil {
			return ErrCommandFailed.ToRESP()
		}

		if wasDeleted {
			deleted++
		}
	}

	return resp.Integer(deleted).ToRESP()
}

package command

import (
	"strings"
	"unicode/utf8"

	"github.com/isaialcantara/toyredis/internal/resp"
)

type commandHandler func(resp.BulkArray) []byte

type commandNode struct {
	handler  commandHandler
	children map[string]*commandNode
}

type CommandDispatcher struct {
	rootCommand *commandNode
}

func NewCommandDispatcher() *CommandDispatcher {
	dispatcher := &CommandDispatcher{
		rootCommand: &commandNode{
			children: make(map[string]*commandNode),
		},
	}

	dispatcher.
		registerCommand([]string{"PING"}, pingHandler).
		registerCommand([]string{"ECHO"}, echoHandler)

	return dispatcher
}

func (d *CommandDispatcher) Dispatch(bulkArray resp.BulkArray) []byte {
	if len(bulkArray) < 1 {
		return ErrCommandEmpty.ToRESP()
	}

	node := d.rootCommand
	consumed := 0
	for _, bulkString := range bulkArray {
		if !utf8.Valid(bulkString) {
			return ErrCommandInvalid.ToRESP()
		}
		key := strings.ToUpper(string(bulkString))
		if nextCommand, ok := node.children[key]; ok {
			node = nextCommand
			consumed++
		} else {
			break
		}
	}

	if node.handler != nil {
		return node.handler(bulkArray[consumed:])
	}
	return ErrCommandInvalid.ToRESP()
}

func (d *CommandDispatcher) registerCommand(path []string, handler commandHandler) *CommandDispatcher {
	node := d.rootCommand
	for _, commandPart := range path {
		key := strings.ToUpper(commandPart)
		if node.children[key] == nil {
			node.children[key] = &commandNode{children: make(map[string]*commandNode)}
		}
		node = node.children[key]
	}
	node.handler = handler
	return d
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
	return resp.SimpleString(string(args[0])).ToRESP()
}

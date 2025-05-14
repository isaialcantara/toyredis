package command

import (
	"strings"
	"unicode/utf8"

	"github.com/isaialcantara/toyredis/internal/resp"
)

var rootCommand *commandNode = &commandNode{
	children: make(map[string]*commandNode),
}

type commandHandler func(resp.BulkArray) (resp.RESPType, error)

type commandNode struct {
	handler  commandHandler
	children map[string]*commandNode
}

func init() {
	registerCommand([]string{"PING"}, pingHandler)
	registerCommand([]string{"ECHO"}, echoHandler)
}

func registerCommand(path []string, handler commandHandler) {
	node := rootCommand
	for _, commandPart := range path {
		key := strings.ToUpper(commandPart)
		if node.children[key] == nil {
			node.children[key] = &commandNode{children: make(map[string]*commandNode)}
		}
		node = node.children[key]
	}
	node.handler = handler
}

func DispatchCommand(bulkArray resp.BulkArray) (resp.RESPType, error) {
	if len(bulkArray) < 1 {
		return nil, ErrCommandEmpty
	}

	node := rootCommand
	consumed := 0
	for _, bulkString := range bulkArray {
		if !utf8.Valid(bulkString) {
			return nil, ErrCommandInvalid
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
	return nil, ErrCommandInvalid
}

func pingHandler(args resp.BulkArray) (resp.RESPType, error) {
	if len(args) == 0 {
		return resp.SimpleString("PONG"), nil
	}
	return echoHandler(args)
}

func echoHandler(args resp.BulkArray) (resp.RESPType, error) {
	if len(args) != 1 {
		return nil, ErrCommandArgsNumber
	}
	return resp.SimpleString(string(args[0])), nil
}

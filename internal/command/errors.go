package command

import "github.com/isaialcantara/toyredis/internal/resp"

const (
	ErrCommandEmpty      = resp.SimpleError("ERR: empty command")
	ErrCommandInvalid    = resp.SimpleError("ERR: invalid command")
	ErrCommandArgsNumber = resp.SimpleError("ERR: wrong number or arguments")
	ErrCommandFailed     = resp.SimpleError("ERR: command failed")
)

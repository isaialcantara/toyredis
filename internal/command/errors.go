package command

import "github.com/isaialcantara/toyredis/internal/resp"

const (
	ErrCommandEmpty      = resp.SimpleError("empty command")
	ErrCommandInvalid    = resp.SimpleError("invalid command")
	ErrCommandArgsNumber = resp.SimpleError("wrong number or arguments")
)

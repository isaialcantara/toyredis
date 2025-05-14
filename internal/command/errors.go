package command

type CommandError string

const (
	ErrCommandEmpty      = CommandError("empty command")
	ErrCommandInvalid    = CommandError("invalid command")
	ErrCommandArgsNumber = CommandError("wrong number or arguments")
)

func (e CommandError) Error() string {
	return "ERR Command error: " + string(e)
}

func (e CommandError) ToRESP() []byte {
	withPrefix := "-" + e.Error()
	return append([]byte(withPrefix), '\r', '\n')
}

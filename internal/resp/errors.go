package resp

type ProtocolError string

func (e ProtocolError) Error() string {
	return "ERR Protocol error: " + string(e)
}

func (e ProtocolError) ToRESP() []byte {
	withPrefix := append([]byte{'-'}, []byte(e.Error())...)
	return append(withPrefix, '\r', '\n')
}

const (
	// Tokenizer
	ErrProtocolInvalidType          = ProtocolError("invalid input type")
	ErrProtocolInvalidBulkArrLength = ProtocolError("invalid bulk string array length")
	ErrProtocolInvalidBulkLength    = ProtocolError("invalid bulk string length")
	ErrProtocolNoCRLF               = ProtocolError("line was not terminated with a CRLF")
	ErrProtocolMissingBulkData      = ProtocolError("missing bulk string data")

	// Parser
	ErrProtocolNotBulkArray         = ProtocolError("input isn't a bulk string array")
	ErrProtocolIncompleteBulkArray  = ProtocolError("input bulk array is incomplete")
	ErrProtocolIncompleteBulkString = ProtocolError("input bulk string is incomplete")
)

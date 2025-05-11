package resp

type customError string

func (e customError) Error() string {
	return string(e)
}

const (
	ErrProtocolInvalidType          = customError("ERR Protocol error: invalid input type")
	ErrProtocolInvalidBulkArrLength = customError("ERR Protocol error: invalid bulk string array length")
	ErrProtocolInvalidBulkLength    = customError("ERR Protocol error: invalid bulk string length")
	ErrProtocolNoCRLF               = customError("ERR Protocol error: line was not terminated with a CRLF")
	ErrProtocolMissingBulkData      = customError("ERR Protocol error: missing bulk string data")
)

package resp

type RESPError interface {
	error
	RESPType
}

type ProtocolError string

func (e ProtocolError) Error() string {
	return "ERR Protocol error: " + string(e)
}

func (e ProtocolError) ToRESP() []byte {
	return []byte("-" + e.Error() + "\r\n")
}

const (
	ErrProtocolNoCRLF                 = ProtocolError("line was not terminated with a CRLF")
	ErrProtocolInvalidBulkArrayLength = ProtocolError("invalid bulk string array length")
	ErrProtocolInvalidBulkLength      = ProtocolError("invalid bulk string length")
	ErrProtocolMissingBulkData        = ProtocolError("missing bulk string data")
	ErrProtocolNoBulkStart            = ProtocolError("no '$' was sent for bulk start")
	ErrProtocolUnbalancedQuotes       = ProtocolError("unbalanced quotes in request")
)

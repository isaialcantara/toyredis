package resp

type Tokenizer interface {
	NextToken() (Token, error)
}

type Token struct {
	Type   string
	Data   []byte
	Length int64
}

func newBulkArrayStartToken(length int64) Token  { return Token{Type: "bulkArrayStart", Length: length} }
func newBulkStringStartToken(length int64) Token { return Token{Type: "bulkStart", Length: length} }
func newBulkDataToken(data []byte) Token         { return Token{Type: "bulkData", Data: data} }

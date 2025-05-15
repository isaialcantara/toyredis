package resp

type Tokenizer interface {
	NextToken() (Token, error)
}

type Token struct {
	Type   string
	Length int64
	Data   []byte
}

func newBulkArrayStartToken(length int64) Token  { return Token{Type: "bulkArrayStart", Length: length} }
func newBulkStringStartToken(length int64) Token { return Token{Type: "bulkStart", Length: length} }
func newBulkDataToken(data []byte) Token         { return Token{Type: "bulkData", Data: data} }

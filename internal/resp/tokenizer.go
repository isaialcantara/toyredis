package resp

type Tokenizer interface {
	NextToken() (Token, error)
}

type Token any

type (
	BulkArrayStartToken  struct{ Length int64 }
	BulkStringStartToken struct{ Length int64 }
	BulkDataToken        struct{ Data []byte }
)

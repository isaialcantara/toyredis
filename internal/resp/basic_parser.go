package resp

type BasicParser struct {
	tokenizer Tokenizer
}

func NewBasicParser(tokenizer Tokenizer) *BasicParser {
	return &BasicParser{tokenizer: tokenizer}
}

func (p *BasicParser) NextBulkArray() (BulkArray, error) {
	token, err := p.tokenizer.NextToken()
	if err != nil {
		return BulkArray{}, err
	}

	arrStartToken, ok := token.(BulkArrayStartToken)
	if !ok {
		return BulkArray{}, ErrProtocolNotBulkArray
	}

	bulkArray := BulkArray{}

	if arrStartToken.Length < 1 {
		return bulkArray, nil
	}

	return p.readBulkArrayItems(bulkArray, arrStartToken.Length)
}

func (p *BasicParser) readBulkArrayItems(bulkArray BulkArray, arrayLength int64) (BulkArray, error) {
	for range arrayLength {
		startToken, err := p.tokenizer.NextToken()
		if err != nil {
			return BulkArray{}, err
		}

		bulkStringStartToken, ok := startToken.(BulkStringStartToken)
		if !ok {
			return BulkArray{}, ErrProtocolIncompleteBulkArray
		}

		bulkString := BulkString{}

		if bulkStringStartToken.Length > 0 {
			token, err := p.tokenizer.NextToken()
			if err != nil {
				return BulkArray{}, err
			}

			if bulkDataToken, ok := token.(BulkDataToken); ok {
				bulkString = bulkDataToken.Data
			} else {
				return BulkArray{}, ErrProtocolIncompleteBulkString
			}
		}

		if bulkStringStartToken.Length != int64(len(bulkString)) {
			return BulkArray{}, ErrProtocolIncompleteBulkString
		}

		bulkArray = append(bulkArray, bulkString)
	}

	return bulkArray, nil
}

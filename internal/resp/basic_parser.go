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

	bulkArray := BulkArray{declaredLength: arrStartToken.Length, BulkStrings: []BulkString{}}
	if bulkArray.declaredLength < 1 {
		return bulkArray, nil
	}

	return p.readBulkArrayItems(bulkArray)
}

func (p *BasicParser) readBulkArrayItems(bulkArray BulkArray) (BulkArray, error) {
	for range bulkArray.declaredLength {
		startToken, err := p.tokenizer.NextToken()
		if err != nil {
			return BulkArray{}, err
		}

		bulkStringStartToken, ok := startToken.(BulkStringStartToken)
		if !ok {
			return BulkArray{}, ErrProtocolIncompleteBulkArray
		}

		bulkString := BulkString{declaredLength: bulkStringStartToken.Length}

		if bulkString.declaredLength == 0 {
			bulkString.Data = []byte{}
		}

		if bulkString.declaredLength > 0 {
			token, err := p.tokenizer.NextToken()
			if err != nil {
				return BulkArray{}, err
			}

			if bulkData, ok := token.(BulkDataToken); ok {
				bulkString.Data = bulkData.Data
			} else {
				return BulkArray{}, ErrProtocolIncompleteBulkString
			}
		}

		bulkArray.BulkStrings = append(bulkArray.BulkStrings, bulkString)
	}

	return bulkArray, nil
}

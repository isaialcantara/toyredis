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

	if token.Type != "bulkArrayStart" {
		return BulkArray{}, ErrProtocolNotBulkArray
	}

	bulkArray := BulkArray{}

	if token.Length < 1 {
		return bulkArray, nil
	}

	return p.readBulkArrayItems(bulkArray, token.Length)
}

func (p *BasicParser) readBulkArrayItems(bulkArray BulkArray, arrayLength int64) (BulkArray, error) {
	for range arrayLength {
		startToken, err := p.tokenizer.NextToken()
		if err != nil {
			return BulkArray{}, err
		}

		if startToken.Type != "bulkStart" {
			return BulkArray{}, ErrProtocolIncompleteBulkArray
		}

		bulkString := BulkString{}

		if startToken.Length > 0 {
			dataToken, err := p.tokenizer.NextToken()
			if err != nil {
				return BulkArray{}, err
			}

			if dataToken.Type == "bulkData" {
				bulkString = dataToken.Data
			} else {
				return BulkArray{}, ErrProtocolIncompleteBulkString
			}
		}

		if startToken.Length != int64(len(bulkString)) {
			return BulkArray{}, ErrProtocolIncompleteBulkString
		}

		bulkArray = append(bulkArray, bulkString)
	}

	return bulkArray, nil
}

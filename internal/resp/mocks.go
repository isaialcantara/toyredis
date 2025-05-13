package resp

type MockTokenizer struct {
	tokens []Token
	err    error
}

func (t *MockTokenizer) NextToken() (Token, error) {
	if len(t.tokens) == 0 {
		return nil, t.err
	}

	token := t.tokens[0]
	t.tokens = t.tokens[1:]
	return token, nil
}

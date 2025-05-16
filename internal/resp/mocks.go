package resp

type MockTokenizer struct {
	err    error
	tokens []Token
}

func (t *MockTokenizer) NextToken() (Token, error) {
	if len(t.tokens) == 0 {
		return Token{}, t.err
	}

	token := t.tokens[0]
	t.tokens = t.tokens[1:]
	return token, nil
}

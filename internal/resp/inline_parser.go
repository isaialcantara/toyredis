package resp

import (
	"bufio"
	"bytes"
	"io"
)

type InlineParser struct {
	reader *bufio.Reader
}

var _ Parser = (*InlineParser)(nil)

func NewInlineParser(r io.Reader) *InlineParser {
	return &InlineParser{bufio.NewReader(r)}
}

func (p *InlineParser) NextBulkArray() (BulkArray, error) {
	line, err := p.readLine()
	if err != nil {
		return nil, err
	}

	return splitInlineBulkString(line)
}

func (p *InlineParser) readLine() ([]byte, error) {
	line, err := p.reader.ReadBytes('\n')
	if err != nil {
		return nil, err
	}

	if !bytes.HasSuffix(line, []byte("\r\n")) {
		return nil, ErrProtocolNoCRLF
	}

	return bytes.TrimSuffix(line, []byte("\r\n")), nil
}

func splitInlineBulkString(line []byte) (BulkArray, error) {
	result := BulkArray{}
	token := BulkString{}
	inQuotes := false
	escaped := false

	for _, ch := range line {
		switch {
		case escaped:
			token = append(token, ch)
			escaped = false

		case ch == '\\':
			escaped = true

		case ch == '"':
			inQuotes = !inQuotes

		case ch == ' ' && !inQuotes:
			if len(token) > 0 {
				result = append(result, token)
				token = BulkString{}
			}

		default:
			token = append(token, ch)
		}
	}

	if inQuotes {
		return nil, ErrProtocolUnbalancedQuotes
	}

	if len(token) > 0 {
		result = append(result, token)
	}

	return result, nil
}

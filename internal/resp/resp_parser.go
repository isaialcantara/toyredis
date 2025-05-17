package resp

import (
	"bufio"
	"bytes"
	"io"
	"strconv"
	"strings"
)

type RESPParser struct {
	reader *bufio.Reader
}

var _ Parser = (*RESPParser)(nil)

func NewRESPParser(r io.Reader) *RESPParser {
	return &RESPParser{reader: bufio.NewReader(r)}
}

func (p *RESPParser) NextBulkArray() (BulkArray, error) {
	b, err := p.reader.Peek(1)
	if err != nil {
		return nil, err
	}

	if b[0] == '*' {
		return p.readBulkArray()
	}
	return p.readInlineBulkArray()
}

func (p *RESPParser) readBulkArray() (BulkArray, error) {
	str, err := p.reader.ReadString('\n')
	if err != nil {
		return nil, err
	}

	if !strings.HasSuffix(str, "\r\n") {
		return nil, ErrProtocolNoCRLF
	}

	str = strings.TrimPrefix(str, "*")
	str = strings.TrimSuffix(str, "\r\n")

	length, err := strconv.ParseInt(str, 10, 64)
	if err != nil || length < 0 {
		return nil, ErrProtocolInvalidBulkArrayLength
	}

	bulkArray := make(BulkArray, length)

	for i := range length {
		bulk, err := p.readBulk()
		if err != nil {
			return nil, err
		}

		bulkArray[i] = bulk
	}

	return bulkArray, nil
}

func (p *RESPParser) readBulk() (BulkString, error) {
	str, err := p.reader.ReadString('\n')
	if err != nil {
		return nil, err
	}

	if !strings.HasSuffix(str, "\r\n") {
		return nil, ErrProtocolNoCRLF
	}

	if !strings.HasPrefix(str, "$") {
		return nil, ErrProtocolNoBulkStart
	}

	str = strings.TrimPrefix(str, "$")
	str = strings.TrimSuffix(str, "\r\n")

	length, err := strconv.ParseInt(str, 10, 64)
	if err != nil || length < 0 {
		return nil, ErrProtocolInvalidBulkLength
	}

	buf := make([]byte, length+2)
	_, err = io.ReadFull(p.reader, buf)
	if err != nil {
		return nil, ErrProtocolMissingBulkData
	}

	if !bytes.HasSuffix(buf, []byte("\r\n")) {
		return nil, ErrProtocolNoCRLF
	}

	buf = bytes.TrimSuffix(buf, []byte("\r\n"))

	return BulkString(buf), nil
}

func (p *RESPParser) readInlineBulkArray() (BulkArray, error) {
	inlineParser := NewInlineParser(p.reader)
	return inlineParser.NextBulkArray()
}

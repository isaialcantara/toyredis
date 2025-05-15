package resp

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"strconv"
	"strings"
)

type FSMTokenizer struct {
	reader            *bufio.Reader
	state             tokenizerState
	currentBulkLength int64
}

type tokenizerState int

const (
	stateType tokenizerState = iota
	stateArrayStart
	stateBulkStart
	stateBulkData
)

func NewFSMTokenizer(r io.Reader) *FSMTokenizer {
	return &FSMTokenizer{
		reader: bufio.NewReader(r),
		state:  stateType,
	}
}

func (t *FSMTokenizer) NextToken() (Token, error) {
	switch t.state {
	case stateType:
		return t.handleType()

	case stateArrayStart:
		return t.handleArrayStart()

	case stateBulkStart:
		return t.handleBulkStart()

	case stateBulkData:
		return t.handleBulkData()

	default:
		log.Panicf("Invalid tokenizer state %+v\n", t.state)
		return Token{}, nil
	}
}

func (t *FSMTokenizer) handleType() (Token, error) {
	b, err := t.reader.ReadByte()
	if err != nil {
		return Token{}, err
	}

	switch b {
	case '*':
		t.state = stateArrayStart
	case '$':
		t.state = stateBulkStart

	default:
		return Token{}, ErrProtocolInvalidType
	}

	return t.NextToken()
}

func (t *FSMTokenizer) handleArrayStart() (Token, error) {
	str, err := t.reader.ReadString('\n')
	if err != nil {
		return Token{}, ErrProtocolInvalidBulkArrLength
	}

	if !strings.HasSuffix(str, "\r\n") {
		return Token{}, ErrProtocolNoCRLF
	}

	str = strings.TrimSuffix(str, "\r\n")

	length, err := strconv.ParseInt(str, 10, 64)
	if err != nil || length < 0 {
		return Token{}, ErrProtocolInvalidBulkArrLength
	}

	t.state = stateType
	return newBulkArrayStartToken(length), nil
}

func (t *FSMTokenizer) handleBulkStart() (Token, error) {
	str, err := t.reader.ReadString('\n')
	if err != nil {
		return Token{}, ErrProtocolInvalidBulkLength
	}

	if !strings.HasSuffix(str, "\r\n") {
		return Token{}, ErrProtocolNoCRLF
	}

	str = strings.TrimSuffix(str, "\r\n")

	length, err := strconv.ParseInt(str, 10, 64)
	if err != nil || length < 0 {
		return Token{}, ErrProtocolInvalidBulkLength
	}

	if length > 0 {
		t.currentBulkLength = length
		t.state = stateBulkData
	} else {
		t.state = stateType
	}

	return newBulkStringStartToken(length), nil
}

func (t *FSMTokenizer) handleBulkData() (Token, error) {
	buf := make([]byte, t.currentBulkLength+2)
	_, err := io.ReadFull(t.reader, buf)
	if err != nil {
		return Token{}, ErrProtocolMissingBulkData
	}

	if !bytes.HasSuffix(buf, []byte("\r\n")) {
		return Token{}, ErrProtocolNoCRLF
	}

	buf = bytes.TrimSuffix(buf, []byte("\r\n"))

	t.state = stateType
	return newBulkDataToken(buf), nil
}

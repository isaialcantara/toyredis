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
	stateArrStart
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

	case stateArrStart:
		return t.handleArrStart()

	case stateBulkStart:
		return t.handleBulkStart()

	case stateBulkData:
		return t.handleBulkData()

	default:
		log.Panicf("Invalid tokenizer state %+v\n", t.state)
		return nil, nil
	}
}

func (t *FSMTokenizer) handleType() (Token, error) {
	byte, err := t.reader.ReadByte()
	if err != nil {
		return nil, err
	}

	switch byte {
	case '*':
		t.state = stateArrStart
	case '$':
		t.state = stateBulkStart

	default:
		return nil, ErrProtocolInvalidType
	}

	return t.NextToken()
}

func (t *FSMTokenizer) handleArrStart() (Token, error) {
	str, err := t.reader.ReadString('\n')
	if err != nil {
		return nil, ErrProtocolInvalidBulkArrLength
	}

	if !strings.HasSuffix(str, "\r\n") {
		return nil, ErrProtocolNoCRLF
	}

	str = strings.TrimSuffix(str, "\r\n")

	length, err := strconv.ParseInt(str, 10, 64)
	if err != nil || length < 0 {
		return nil, ErrProtocolInvalidBulkArrLength
	}

	t.state = stateType
	return BulkArrayStartToken{Length: length}, nil
}

func (t *FSMTokenizer) handleBulkStart() (Token, error) {
	str, err := t.reader.ReadString('\n')
	if err != nil {
		return nil, ErrProtocolInvalidBulkLength
	}

	if !strings.HasSuffix(str, "\r\n") {
		return nil, ErrProtocolNoCRLF
	}

	str = strings.TrimSuffix(str, "\r\n")

	length, err := strconv.ParseInt(str, 10, 64)
	if err != nil || length < 0 {
		return nil, ErrProtocolInvalidBulkLength
	}

	if length > 0 {
		t.currentBulkLength = length
		t.state = stateBulkData
	} else {
		t.state = stateType
	}

	return BulkStringStartToken{Length: length}, nil
}

func (t *FSMTokenizer) handleBulkData() (Token, error) {
	buf := make([]byte, t.currentBulkLength+2)
	_, err := io.ReadFull(t.reader, buf)
	if err != nil {
		return nil, ErrProtocolMissingBulkData
	}

	if !bytes.HasSuffix(buf, []byte("\r\n")) {
		return nil, ErrProtocolNoCRLF
	}

	buf = bytes.TrimSuffix(buf, []byte("\r\n"))

	t.state = stateType
	return BulkDataToken{Data: buf}, nil
}

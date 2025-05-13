package resp

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFSMTokenizer_NextToken(t *testing.T) {
	t.Run("return a single token per call", func(t *testing.T) {
		r := strings.NewReader("*3\r\n$3\r\nSET\r\n$3\r\nfoo\r\n$-1\r\n")
		tokenizer := NewFSMTokenizer(r)

		for _, expected := range []Token{
			BulkArrayStartToken{Length: 3},
			BulkStringStartToken{Length: 3},
			BulkDataToken{Data: []byte("SET")},
			BulkStringStartToken{Length: 3},
			BulkDataToken{Data: []byte("foo")},
			BulkStringStartToken{Length: -1},
		} {
			token, err := tokenizer.NextToken()

			assert.NoError(t, err)
			assert.Equal(t, expected, token)
		}
	})

	t.Run("fail with unknown type", func(t *testing.T) {
		r := strings.NewReader("+PING\r\n")
		tokenizer := NewFSMTokenizer(r)

		token, err := tokenizer.NextToken()
		assert.Nil(t, token)
		assert.ErrorIs(t, err, ErrProtocolInvalidType)
	})

	t.Run("return bulk array start token", func(t *testing.T) {
		for _, tc := range []struct {
			line  string
			token Token
			err   error
		}{
			{"*1\r\n", BulkArrayStartToken{Length: 1}, nil},
			{"*0\r\n", BulkArrayStartToken{Length: 0}, nil},
			{"*-1\r\n", nil, ErrProtocolInvalidBulkArrLength},
			{"*abc\r\n", nil, ErrProtocolInvalidBulkArrLength},
			{"*123", nil, ErrProtocolInvalidBulkArrLength},
			{"*3.14\r\n", nil, ErrProtocolInvalidBulkArrLength},
			{"*123\n", nil, ErrProtocolNoCRLF},
		} {
			r := strings.NewReader(tc.line)
			tokenizer := NewFSMTokenizer(r)
			token, err := tokenizer.NextToken()
			assert.Equal(t, tc.token, token)
			assert.Equal(t, tc.err, err)
		}
	})

	t.Run("return bulk string start token", func(t *testing.T) {
		for _, tc := range []struct {
			line  string
			token Token
			err   error
		}{
			{"$1\r\n", BulkStringStartToken{Length: 1}, nil},
			{"$0\r\n", BulkStringStartToken{Length: 0}, nil},
			{"$-1\r\n", BulkStringStartToken{Length: -1}, nil},
			{"$-2\r\n", nil, ErrProtocolInvalidBulkLength},
			{"$abc\r\n", nil, ErrProtocolInvalidBulkLength},
			{"$123", nil, ErrProtocolInvalidBulkLength},
			{"$3.14\r\n", nil, ErrProtocolInvalidBulkLength},
			{"$123\n", nil, ErrProtocolNoCRLF},
		} {
			r := strings.NewReader(tc.line)
			tokenizer := NewFSMTokenizer(r)
			token, err := tokenizer.NextToken()
			assert.Equal(t, tc.token, token)
			assert.Equal(t, tc.err, err)
		}
	})

	t.Run("return bulk data token", func(t *testing.T) {
		for _, tc := range []struct {
			line  string
			token Token
			err   error
		}{
			{"$4\r\nPING\r\n", BulkDataToken{Data: []byte("PING")}, nil},
			{"$4\r\nP\r\nG\r\n", BulkDataToken{Data: []byte{'P', '\r', '\n', 'G'}}, nil},
			{"$4\r\nP\xFFNG\r\n", BulkDataToken{Data: []byte{'P', '\xFF', 'N', 'G'}}, nil},
			{"$4\r\nPI\r\n", nil, ErrProtocolMissingBulkData},
			{"$4\r\nPING", nil, ErrProtocolMissingBulkData},
			{"$4\r\nPINGG\n", nil, ErrProtocolNoCRLF},
		} {
			r := strings.NewReader(tc.line)
			tokenizer := NewFSMTokenizer(r)
			_, err := tokenizer.NextToken()
			assert.NoError(t, err)

			token, err := tokenizer.NextToken()
			assert.Equal(t, tc.token, token)
			assert.Equal(t, tc.err, err)
		}
	})

	t.Run("bulk string empty and null don't require bulk data", func(t *testing.T) {
		r := strings.NewReader("*3\r\n$0\r\n$-1\r\n$3\r\nend\r\n")
		tokenizer := NewFSMTokenizer(r)

		for _, expected := range []Token{
			BulkArrayStartToken{Length: 3},
			BulkStringStartToken{Length: 0},
			BulkStringStartToken{Length: -1},
			BulkStringStartToken{Length: 3},
			BulkDataToken{Data: []byte("end")},
		} {
			token, err := tokenizer.NextToken()
			assert.NoError(t, err)
			assert.Equal(t, expected, token)
		}
	})
}

package resp

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFSMTokenizer_NextToken(t *testing.T) {
	t.Run("return a single token per call", func(t *testing.T) {
		r := strings.NewReader("*3\r\n$3\r\nSET\r\n$3\r\nfoo\r\n")
		tokenizer := NewFSMTokenizer(r)

		for _, expected := range []Token{
			{Type: "bulkArrayStart", Length: 3},
			{Type: "bulkStart", Length: 3},
			{Type: "bulkData", Data: []byte("SET")},
			{Type: "bulkStart", Length: 3},
			{Type: "bulkData", Data: []byte("foo")},
		} {
			token, err := tokenizer.NextToken()

			require.NoError(t, err)
			require.Equal(t, expected, token)
		}
	})

	t.Run("fail with unknown type", func(t *testing.T) {
		r := strings.NewReader("+PING\r\n")
		tokenizer := NewFSMTokenizer(r)

		_, err := tokenizer.NextToken()
		require.ErrorIs(t, err, ErrProtocolInvalidType)
	})

	t.Run("return bulk array start token", func(t *testing.T) {
		for _, tc := range []struct {
			line  string
			token Token
			err   error
		}{
			{"*1\r\n", Token{Type: "bulkArrayStart", Length: 1}, nil},
			{"*0\r\n", Token{Type: "bulkArrayStart", Length: 0}, nil},
			{"*-1\r\n", Token{}, ErrProtocolInvalidBulkArrLength},
			{"*abc\r\n", Token{}, ErrProtocolInvalidBulkArrLength},
			{"*123", Token{}, ErrProtocolInvalidBulkArrLength},
			{"*3.14\r\n", Token{}, ErrProtocolInvalidBulkArrLength},
			{"*123\n", Token{}, ErrProtocolNoCRLF},
		} {
			r := strings.NewReader(tc.line)
			tokenizer := NewFSMTokenizer(r)
			token, err := tokenizer.NextToken()
			require.Equal(t, tc.token, token)
			require.Equal(t, tc.err, err)
		}
	})

	t.Run("return bulk string start token", func(t *testing.T) {
		for _, tc := range []struct {
			line  string
			token Token
			err   error
		}{
			{"$1\r\n", Token{Type: "bulkStart", Length: 1}, nil},
			{"$0\r\n", Token{Type: "bulkStart", Length: 0}, nil},
			{"$-1\r\n", Token{}, ErrProtocolInvalidBulkLength},
			{"$-2\r\n", Token{}, ErrProtocolInvalidBulkLength},
			{"$abc\r\n", Token{}, ErrProtocolInvalidBulkLength},
			{"$123", Token{}, ErrProtocolInvalidBulkLength},
			{"$3.14\r\n", Token{}, ErrProtocolInvalidBulkLength},
			{"$123\n", Token{}, ErrProtocolNoCRLF},
		} {
			r := strings.NewReader(tc.line)
			tokenizer := NewFSMTokenizer(r)
			token, err := tokenizer.NextToken()
			require.Equal(t, tc.token, token)
			require.Equal(t, tc.err, err)
		}
	})

	t.Run("return bulk data token", func(t *testing.T) {
		for _, tc := range []struct {
			line  string
			token Token
			err   error
		}{
			{"$4\r\nPING\r\n", Token{Type: "bulkData", Data: []byte("PING")}, nil},
			{"$4\r\nP\r\nG\r\n", Token{Type: "bulkData", Data: []byte{'P', '\r', '\n', 'G'}}, nil},
			{"$4\r\nP\xFFNG\r\n", Token{Type: "bulkData", Data: []byte{'P', '\xFF', 'N', 'G'}}, nil},
			{"$4\r\nPI\r\n", Token{}, ErrProtocolMissingBulkData},
			{"$4\r\nPING", Token{}, ErrProtocolMissingBulkData},
			{"$4\r\nPINGG\n", Token{}, ErrProtocolNoCRLF},
		} {
			r := strings.NewReader(tc.line)
			tokenizer := NewFSMTokenizer(r)
			_, err := tokenizer.NextToken()
			require.NoError(t, err)

			token, err := tokenizer.NextToken()
			require.Equal(t, tc.token, token)
			require.Equal(t, tc.err, err)
		}
	})

	t.Run("empty bulk string doesn't require bulk data", func(t *testing.T) {
		r := strings.NewReader("*2\r\n$0\r\n$3\r\nend\r\n")
		tokenizer := NewFSMTokenizer(r)

		for _, expected := range []Token{
			{Type: "bulkArrayStart", Length: 2},
			{Type: "bulkStart", Length: 0},
			{Type: "bulkStart", Length: 3},
			{Type: "bulkData", Data: []byte("end")},
		} {
			token, err := tokenizer.NextToken()
			require.NoError(t, err)
			require.Equal(t, expected, token)
		}
	})
}

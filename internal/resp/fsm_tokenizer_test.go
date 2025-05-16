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
			err   error
			token Token
		}{
			{"*1\r\n", nil, Token{Type: "bulkArrayStart", Length: 1}},
			{"*0\r\n", nil, Token{Type: "bulkArrayStart", Length: 0}},
			{"*-1\r\n", ErrProtocolInvalidBulkArrLength, Token{}},
			{"*abc\r\n", ErrProtocolInvalidBulkArrLength, Token{}},
			{"*123", ErrProtocolInvalidBulkArrLength, Token{}},
			{"*3.14\r\n", ErrProtocolInvalidBulkArrLength, Token{}},
			{"*123\n", ErrProtocolNoCRLF, Token{}},
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
			err   error
			token Token
		}{
			{"$1\r\n", nil, Token{Type: "bulkStart", Length: 1}},
			{"$0\r\n", nil, Token{Type: "bulkStart", Length: 0}},
			{"$-1\r\n", ErrProtocolInvalidBulkLength, Token{}},
			{"$-2\r\n", ErrProtocolInvalidBulkLength, Token{}},
			{"$abc\r\n", ErrProtocolInvalidBulkLength, Token{}},
			{"$123", ErrProtocolInvalidBulkLength, Token{}},
			{"$3.14\r\n", ErrProtocolInvalidBulkLength, Token{}},
			{"$123\n", ErrProtocolNoCRLF, Token{}},
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
			err   error
			token Token
		}{
			{"$4\r\nPING\r\n", nil, Token{Type: "bulkData", Data: []byte("PING")}},
			{"$4\r\nP\r\nG\r\n", nil, Token{Type: "bulkData", Data: []byte{'P', '\r', '\n', 'G'}}},
			{"$4\r\nP\xFFNG\r\n", nil, Token{Type: "bulkData", Data: []byte{'P', '\xFF', 'N', 'G'}}},
			{"$4\r\nPI\r\n", ErrProtocolMissingBulkData, Token{}},
			{"$4\r\nPING", ErrProtocolMissingBulkData, Token{}},
			{"$4\r\nPINGG\n", ErrProtocolNoCRLF, Token{}},
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

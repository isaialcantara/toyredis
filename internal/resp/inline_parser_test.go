package resp

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInlineParser_NextBulkArray(t *testing.T) {
	t.Run("return bulk array with inline command", func(t *testing.T) {
		r := strings.NewReader(`ECHO abc "def ghi" jkl "mno pqrs" tuv "wxyz"` + "\r\n")
		parser := NewInlineParser(r)

		expected := BulkArray{
			BulkString("ECHO"),
			BulkString("abc"),
			BulkString("def ghi"),
			BulkString("jkl"),
			BulkString("mno pqrs"),
			BulkString("tuv"),
			BulkString("wxyz"),
		}

		bulkArray, err := parser.NextBulkArray()
		require.NoError(t, err)
		require.Equal(t, expected, bulkArray)
	})

	t.Run("return error when the next line does not end with CRLF", func(t *testing.T) {
		r := strings.NewReader("echo abc\n")
		parser := NewInlineParser(r)
		_, err := parser.NextBulkArray()
		require.ErrorIs(t, err, ErrProtocolNoCRLF)
	})

	t.Run("return error with unbalanced quotes", func(t *testing.T) {
		r := strings.NewReader("echo \"abc\r\n")
		parser := NewInlineParser(r)
		_, err := parser.NextBulkArray()
		require.ErrorIs(t, err, ErrProtocolUnbalancedQuotes)
	})
}

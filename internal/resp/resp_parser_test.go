package resp

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRESPParser_NextBulkArray(t *testing.T) {
	t.Run("return bulk array when command is multibulk", func(t *testing.T) {
		r := strings.NewReader("*2\r\n$4\r\nECHO\r\n$12\r\nHello World!\r\n")
		parser := NewRESPParser(r)

		expected := BulkArray{
			BulkString("ECHO"),
			BulkString("Hello World!"),
		}

		bulkArray, err := parser.NextBulkArray()
		require.NoError(t, err)
		require.Equal(t, expected, bulkArray)
	})

	t.Run("return bulk array when command is inline", func(t *testing.T) {
		r := strings.NewReader(`ECHO abc "def ghi" jkl "mno pqrs" tuv "wxyz"` + "\r\n")
		parser := NewRESPParser(r)

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
}

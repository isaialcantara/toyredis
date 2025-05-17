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

	t.Run("fallback to inline parser when command is inline", func(t *testing.T) {
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

	t.Run("return error with no CRLF terminator for bulk array length", func(t *testing.T) {
		r := strings.NewReader("*1\n")
		parser := NewRESPParser(r)

		_, err := parser.NextBulkArray()
		require.ErrorIs(t, err, ErrProtocolNoCRLF)
	})

	t.Run("return error with non integer bulk array length", func(t *testing.T) {
		r := strings.NewReader("*abc\r\n")
		parser := NewRESPParser(r)

		_, err := parser.NextBulkArray()
		require.ErrorIs(t, err, ErrProtocolInvalidBulkArrayLength)
	})

	t.Run("return error with negative bulk array length", func(t *testing.T) {
		r := strings.NewReader("*-1\r\n")
		parser := NewRESPParser(r)

		_, err := parser.NextBulkArray()
		require.ErrorIs(t, err, ErrProtocolInvalidBulkArrayLength)
	})

	t.Run("return error with no CRLF terminator for bulk length", func(t *testing.T) {
		r := strings.NewReader("*1\r\n$3\n")
		parser := NewRESPParser(r)

		_, err := parser.NextBulkArray()
		require.ErrorIs(t, err, ErrProtocolNoCRLF)
	})

	t.Run("return error with invalid bulk start char", func(t *testing.T) {
		r := strings.NewReader("*1\r\n+4\r\nPING\r\n")
		parser := NewRESPParser(r)

		_, err := parser.NextBulkArray()
		require.ErrorIs(t, err, ErrProtocolNoBulkStart)
	})

	t.Run("return error with missing bulk start char", func(t *testing.T) {
		r := strings.NewReader("*1\r\n4\r\nPING\r\n")
		parser := NewRESPParser(r)

		_, err := parser.NextBulkArray()
		require.ErrorIs(t, err, ErrProtocolNoBulkStart)
	})

	t.Run("return error with non integer bulk length", func(t *testing.T) {
		r := strings.NewReader("*1\r\n$abc\r\n")
		parser := NewRESPParser(r)

		_, err := parser.NextBulkArray()
		require.ErrorIs(t, err, ErrProtocolInvalidBulkLength)
	})

	t.Run("return error with negative bulk length", func(t *testing.T) {
		r := strings.NewReader("*1\r\n$-1\r\n")
		parser := NewRESPParser(r)

		_, err := parser.NextBulkArray()
		require.ErrorIs(t, err, ErrProtocolInvalidBulkLength)
	})

	t.Run("return error with incomplete bulk data", func(t *testing.T) {
		r := strings.NewReader("*1\r\n$4\r\nPIN\r\n")
		parser := NewRESPParser(r)

		_, err := parser.NextBulkArray()
		require.ErrorIs(t, err, ErrProtocolMissingBulkData)
	})

	t.Run("return error with missing bulk data", func(t *testing.T) {
		r := strings.NewReader("*1\r\n$4\r\n")
		parser := NewRESPParser(r)

		_, err := parser.NextBulkArray()
		require.ErrorIs(t, err, ErrProtocolMissingBulkData)
	})

	t.Run("return error with more data than previously declared", func(t *testing.T) {
		r := strings.NewReader("*1\r\n$4\r\nPINGG\r\n")
		parser := NewRESPParser(r)

		_, err := parser.NextBulkArray()
		require.ErrorIs(t, err, ErrProtocolNoCRLF)
	})
}

package resp

import (
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewBasicParser(t *testing.T) {
	t.Run("test 1", func(t *testing.T) {
		r := strings.NewReader("*1\r\n$3\r\nABC\r\n")
		tokenizer := NewFSMTokenizer(r)
		expected := &BasicParser{tokenizer}
		assert.Equal(t, expected, NewBasicParser(tokenizer))
	})
}

func TestNextBulkArray(t *testing.T) {
	t.Run("test SET", func(t *testing.T) {
		r := strings.NewReader("*3\r\n$3\r\nSET\r\n$3\r\nkey\r\n$5\r\nvalue\r\n")
		tokenizer := NewFSMTokenizer(r)
		parser := NewBasicParser(tokenizer)
		expected := BulkArray{
			declaredLength: 3,
			BulkStrings: []BulkString{
				{declaredLength: 3, Data: []byte("SET")},
				{declaredLength: 3, Data: []byte("key")},
				{declaredLength: 5, Data: []byte("value")},
			},
		}

		bulkArray, err := parser.NextBulkArray()
		assert.NoError(t, err)
		assert.Equal(t, expected, bulkArray)
	})

	t.Run("test null and empty bulk strings", func(t *testing.T) {
		r := strings.NewReader("*3\r\n$3\r\nSET\r\n$0\r\n$-1\r\n")
		tokenizer := NewFSMTokenizer(r)
		parser := NewBasicParser(tokenizer)
		expected := BulkArray{
			declaredLength: 3,
			BulkStrings: []BulkString{
				{declaredLength: 3, Data: []byte("SET")},
				{declaredLength: 0, Data: []byte{}},
				{declaredLength: -1, Data: nil},
			},
		}

		bulkArray, err := parser.NextBulkArray()
		assert.NoError(t, err)
		assert.Equal(t, expected, bulkArray)
	})

	t.Run("test empty bulk array", func(t *testing.T) {
		r := strings.NewReader("*0\r\n")
		tokenizer := NewFSMTokenizer(r)
		parser := NewBasicParser(tokenizer)
		expected := BulkArray{declaredLength: 0, BulkStrings: []BulkString{}}

		bulkArray, err := parser.NextBulkArray()
		assert.NoError(t, err)
		assert.Equal(t, expected, bulkArray)
	})

	t.Run("test tokenizer error in array declaration", func(t *testing.T) {
		r := strings.NewReader("*3\n")
		tokenizer := NewFSMTokenizer(r)
		parser := NewBasicParser(tokenizer)

		_, err := parser.NextBulkArray()
		assert.ErrorIs(t, err, ErrProtocolNoCRLF)
	})

	t.Run("test tokenizer error in bulk string declaration", func(t *testing.T) {
		r := strings.NewReader("*3\r\n$abc\r\n")
		tokenizer := NewFSMTokenizer(r)
		parser := NewBasicParser(tokenizer)

		_, err := parser.NextBulkArray()
		assert.ErrorIs(t, err, ErrProtocolInvalidBulkLength)
	})

	t.Run("test tokenizer error in bulk data", func(t *testing.T) {
		r := strings.NewReader("*3\r\n$1\r\na\n")
		tokenizer := NewFSMTokenizer(r)
		parser := NewBasicParser(tokenizer)

		_, err := parser.NextBulkArray()
		assert.ErrorIs(t, err, ErrProtocolMissingBulkData)
	})

	t.Run("test not a bulk array", func(t *testing.T) {
		r := strings.NewReader("$4\r\nPING\r\n")
		tokenizer := NewFSMTokenizer(r)
		parser := NewBasicParser(tokenizer)

		_, err := parser.NextBulkArray()
		assert.ErrorIs(t, err, ErrProtocolNotBulkArray)
	})

	t.Run("test incomplete bulk array", func(t *testing.T) {
		r := strings.NewReader("*2\r\n*0\r\n$4\r\nPING\r\n")
		tokenizer := NewFSMTokenizer(r)
		parser := NewBasicParser(tokenizer)

		_, err := parser.NextBulkArray()
		assert.ErrorIs(t, err, ErrProtocolIncompleteBulkArray)
	})

	t.Run("test EOF before completing bulk array", func(t *testing.T) {
		r := strings.NewReader("*2\r\n$4\r\nPING\r\n")
		tokenizer := NewFSMTokenizer(r)
		parser := NewBasicParser(tokenizer)

		_, err := parser.NextBulkArray()
		assert.ErrorIs(t, err, io.EOF)
	})
}

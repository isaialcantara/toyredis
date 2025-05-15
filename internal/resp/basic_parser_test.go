package resp

import (
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewBasicParser(t *testing.T) {
	t.Run("test 1", func(t *testing.T) {
		tokenizer := &MockTokenizer{tokens: []Token{}, err: nil}
		expected := &BasicParser{tokenizer}
		require.Equal(t, expected, NewBasicParser(tokenizer))
	})
}

func TestBasicParser_NextBulkArray(t *testing.T) {
	t.Run("returns the right RESP types", func(t *testing.T) {
		tokenizer := &MockTokenizer{
			tokens: []Token{
				{Type: "bulkArrayStart", Length: 3},
				{Type: "bulkStart", Length: 3},
				{Type: "bulkData", Data: []byte("SET")},
				{Type: "bulkStart", Length: 3},
				{Type: "bulkData", Data: []byte("key")},
				{Type: "bulkStart", Length: 5},
				{Type: "bulkData", Data: []byte("value")},
			},
		}
		parser := NewBasicParser(tokenizer)
		expected := BulkArray{
			BulkString("SET"),
			BulkString("key"),
			BulkString("value"),
		}

		bulkArray, err := parser.NextBulkArray()
		require.NoError(t, err)
		require.Equal(t, expected, bulkArray)
	})

	t.Run("returns empty bulk strings", func(t *testing.T) {
		tokenizer := &MockTokenizer{
			tokens: []Token{
				{Type: "bulkArrayStart", Length: 2},
				{Type: "bulkStart", Length: 3},
				{Type: "bulkData", Data: []byte("SET")},
				{Type: "bulkStart", Length: 0},
			},
		}
		parser := NewBasicParser(tokenizer)
		expected := BulkArray{
			BulkString("SET"),
			BulkString{},
		}

		bulkArray, err := parser.NextBulkArray()
		require.NoError(t, err)
		require.Equal(t, expected, bulkArray)
	})

	t.Run("returns empty bulk array", func(t *testing.T) {
		tokenizer := &MockTokenizer{
			tokens: []Token{{Type: "bulkArrayStart", Length: 0}},
			err:    nil,
		}
		parser := NewBasicParser(tokenizer)
		expected := BulkArray{}

		bulkArray, err := parser.NextBulkArray()
		require.NoError(t, err)
		require.Equal(t, expected, bulkArray)
	})

	t.Run("returns error with tokenizer error in array declaration", func(t *testing.T) {
		tokenizer := &MockTokenizer{err: ErrProtocolNoCRLF}
		parser := NewBasicParser(tokenizer)

		_, err := parser.NextBulkArray()
		require.ErrorIs(t, err, ErrProtocolNoCRLF)
	})

	t.Run("returns error with tokenizer error in bulk string declaration", func(t *testing.T) {
		tokenizer := &MockTokenizer{
			tokens: []Token{{Type: "bulkArrayStart", Length: 3}},
			err:    ErrProtocolInvalidBulkLength,
		}
		parser := NewBasicParser(tokenizer)

		_, err := parser.NextBulkArray()
		require.ErrorIs(t, err, ErrProtocolInvalidBulkLength)
	})

	t.Run("returns error with tokenizer error in bulk data", func(t *testing.T) {
		tokenizer := &MockTokenizer{
			tokens: []Token{
				{Type: "bulkArrayStart", Length: 3},
				{Type: "bulkStart", Length: 1},
			},
			err: ErrProtocolMissingBulkData,
		}
		parser := NewBasicParser(tokenizer)

		_, err := parser.NextBulkArray()
		require.ErrorIs(t, err, ErrProtocolMissingBulkData)
	})

	t.Run("return error when the input isn't a bulk array", func(t *testing.T) {
		tokenizer := &MockTokenizer{
			tokens: []Token{
				{Type: "bulkStart", Length: 4},
				{Type: "bulkData", Data: []byte("PING")},
			},
		}
		parser := NewBasicParser(tokenizer)

		_, err := parser.NextBulkArray()
		require.ErrorIs(t, err, ErrProtocolNotBulkArray)
	})

	t.Run("returns error with incomplete bulk array", func(t *testing.T) {
		tokenizer := &MockTokenizer{
			tokens: []Token{
				{Type: "bulkArrayStart", Length: 2},
				{Type: "bulkArrayStart", Length: 0},
				{Type: "bulkStart", Length: 4},
				{Type: "bulkData", Data: []byte("PING")},
			},
		}
		parser := NewBasicParser(tokenizer)

		_, err := parser.NextBulkArray()
		require.ErrorIs(t, err, ErrProtocolIncompleteBulkArray)
	})

	t.Run("returns error with bulk string missing it's data field", func(t *testing.T) {
		tokenizer := &MockTokenizer{
			tokens: []Token{
				{Type: "bulkArrayStart", Length: 2},
				{Type: "bulkStart", Length: 4},
				{Type: "bulkStart", Length: 4},
			},
		}
		parser := NewBasicParser(tokenizer)

		_, err := parser.NextBulkArray()
		require.ErrorIs(t, err, ErrProtocolIncompleteBulkString)
	})

	t.Run("returns error with incomplete bulk string", func(t *testing.T) {
		tokenizer := &MockTokenizer{
			tokens: []Token{
				{Type: "bulkArrayStart", Length: 1},
				{Type: "bulkStart", Length: 2},
				{Type: "bulkData", Data: []byte("a")},
			},
		}
		parser := NewBasicParser(tokenizer)

		_, err := parser.NextBulkArray()
		require.ErrorIs(t, err, ErrProtocolIncompleteBulkString)
	})

	t.Run("return error when tokenizer returns EOF before completing bulk array", func(t *testing.T) {
		tokenizer := &MockTokenizer{
			tokens: []Token{
				{Type: "bulkArrayStart", Length: 2},
				{Type: "bulkStart", Length: 4},
				{Type: "bulkData", Data: []byte("PING")},
			},
			err: io.EOF,
		}
		parser := NewBasicParser(tokenizer)

		_, err := parser.NextBulkArray()
		require.ErrorIs(t, err, io.EOF)
	})
}

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
				BulkArrayStartToken{Length: 3},
				BulkStringStartToken{Length: 3},
				BulkDataToken{Data: []byte("SET")},
				BulkStringStartToken{Length: 3},
				BulkDataToken{Data: []byte("key")},
				BulkStringStartToken{Length: 5},
				BulkDataToken{Data: []byte("value")},
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
				BulkArrayStartToken{Length: 2},
				BulkStringStartToken{Length: 3},
				BulkDataToken{Data: []byte("SET")},
				BulkStringStartToken{Length: 0},
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
			tokens: []Token{BulkArrayStartToken{Length: 0}},
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
			tokens: []Token{BulkArrayStartToken{Length: 3}},
			err:    ErrProtocolInvalidBulkLength,
		}
		parser := NewBasicParser(tokenizer)

		_, err := parser.NextBulkArray()
		require.ErrorIs(t, err, ErrProtocolInvalidBulkLength)
	})

	t.Run("returns error with tokenizer error in bulk data", func(t *testing.T) {
		tokenizer := &MockTokenizer{
			tokens: []Token{
				BulkArrayStartToken{Length: 3},
				BulkStringStartToken{Length: 1},
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
				BulkStringStartToken{Length: 4},
				BulkDataToken{Data: []byte("PING")},
			},
		}
		parser := NewBasicParser(tokenizer)

		_, err := parser.NextBulkArray()
		require.ErrorIs(t, err, ErrProtocolNotBulkArray)
	})

	t.Run("returns error with incomplete bulk array", func(t *testing.T) {
		tokenizer := &MockTokenizer{
			tokens: []Token{
				BulkArrayStartToken{Length: 2},
				BulkArrayStartToken{Length: 0},
				BulkStringStartToken{Length: 4},
				BulkDataToken{Data: []byte("PING")},
			},
		}
		parser := NewBasicParser(tokenizer)

		_, err := parser.NextBulkArray()
		require.ErrorIs(t, err, ErrProtocolIncompleteBulkArray)
	})

	t.Run("returns error with bulk string missing it's data field", func(t *testing.T) {
		tokenizer := &MockTokenizer{
			tokens: []Token{
				BulkArrayStartToken{Length: 2},
				BulkStringStartToken{Length: 4},
				BulkStringStartToken{Length: 4},
			},
		}
		parser := NewBasicParser(tokenizer)

		_, err := parser.NextBulkArray()
		require.ErrorIs(t, err, ErrProtocolIncompleteBulkString)
	})

	t.Run("returns error with incomplete bulk string", func(t *testing.T) {
		tokenizer := &MockTokenizer{
			tokens: []Token{
				BulkArrayStartToken{Length: 1},
				BulkStringStartToken{Length: 2},
				BulkDataToken{Data: []byte("a")},
			},
		}
		parser := NewBasicParser(tokenizer)

		_, err := parser.NextBulkArray()
		require.ErrorIs(t, err, ErrProtocolIncompleteBulkString)
	})

	t.Run("return error when tokenizer returns EOF before completing bulk array", func(t *testing.T) {
		tokenizer := &MockTokenizer{
			tokens: []Token{
				BulkArrayStartToken{Length: 2},
				BulkStringStartToken{Length: 4},
				BulkDataToken{Data: []byte("PING")},
			},
			err: io.EOF,
		}
		parser := NewBasicParser(tokenizer)

		_, err := parser.NextBulkArray()
		require.ErrorIs(t, err, io.EOF)
	})
}

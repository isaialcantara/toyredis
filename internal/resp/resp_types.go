package resp

import "fmt"

type RESPType interface {
	ToRESP() []byte
}

type SimpleString string

var _ RESPType = (*SimpleString)(nil)

func (s SimpleString) ToRESP() []byte {
	return []byte("+" + s + "\r\n")
}

type SimpleError string

var _ RESPType = (*SimpleError)(nil)

func (e SimpleError) ToRESP() []byte {
	return []byte("-" + e + "\r\n")
}

type Integer int64

var _ RESPType = (*Integer)(nil)

func (i Integer) ToRESP() []byte {
	str := fmt.Sprintf(":%d\r\n", i)
	return []byte(str)
}

type BulkArray []BulkString

var _ RESPType = (*BulkArray)(nil)

func (a BulkArray) ToRESP() []byte {
	lengthStr := fmt.Sprintf("*%d\r\n", len(a))
	lengthLine := []byte(lengthStr)
	if len(a) < 1 {
		return lengthLine
	}

	bulkStringsLines := []byte{}
	for _, bulkString := range a {
		bulkStringsLines = append(bulkStringsLines, bulkString.ToRESP()...)
	}

	return append(lengthLine, bulkStringsLines...)
}

type BulkString []byte

var _ RESPType = (*BulkString)(nil)

func (s BulkString) ToRESP() []byte {
	lengthStr := fmt.Sprintf("$%d\r\n", len(s))
	lengthLine := []byte(lengthStr)
	if len(s) < 1 {
		return lengthLine
	}

	s = append(s, '\r', '\n')
	return append(lengthLine, s...)
}

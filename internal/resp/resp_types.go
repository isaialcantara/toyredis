package resp

import "fmt"

type RESPType interface {
	ToRESP() []byte
}

type BulkArray []BulkString

func (a BulkArray) ToResp() []byte {
	lengthStr := fmt.Sprintf("$%d\r\n", len(a))
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

func (s BulkString) ToRESP() []byte {
	lengthStr := fmt.Sprintf("$%d\r\n", len(s))
	lengthLine := []byte(lengthStr)
	if len(s) < 1 {
		return lengthLine
	}

	dataLine := append(s, '\r', '\n')
	return append(lengthLine, dataLine...)
}

type Integer int64

func (i Integer) ToRESP() []byte {
	str := fmt.Sprintf(":%d\r\n", i)
	return []byte(str)
}

type SimpleString string

func (s SimpleString) ToRESP() []byte {
	return []byte("+" + s + "\r\n")
}

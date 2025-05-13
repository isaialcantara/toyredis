package resp

import "fmt"

type RESPType interface {
	ToRESP() []byte
}

type BulkArray struct {
	declaredLength int64
	BulkStrings    []BulkString
}

func (a BulkArray) ToResp() []byte {
	lengthStr := fmt.Sprintf("$%d\r\n", a.declaredLength)
	lengthLine := []byte(lengthStr)
	if a.declaredLength < 1 {
		return lengthLine
	}

	bulkStringsLines := []byte{}
	for _, bulkString := range a.BulkStrings {
		bulkStringsLines = append(bulkStringsLines, bulkString.ToRESP()...)
	}

	return append(lengthLine, bulkStringsLines...)
}

type BulkString struct {
	declaredLength int64
	Data           []byte
}

func (s BulkString) ToRESP() []byte {
	lengthStr := fmt.Sprintf("$%d\r\n", s.declaredLength)
	lengthLine := []byte(lengthStr)
	if s.declaredLength < 1 {
		return lengthLine
	}

	dataLine := append(s.Data, '\r', '\n')
	return append(lengthLine, dataLine...)
}

type Integer int64

func (i Integer) ToRESP() []byte {
	str := fmt.Sprintf(":%d\r\n", i)
	return []byte(str)
}

type SimpleString string

func (s SimpleString) ToRESP() []byte {
	withPrefix := append([]byte{'+'}, []byte(s)...)
	return append(withPrefix, '\r', '\n')
}

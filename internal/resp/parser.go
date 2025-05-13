package resp

type Parser interface {
	NextBulkArray() (BulkArray, error)
}

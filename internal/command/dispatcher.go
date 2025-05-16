package command

import "github.com/isaialcantara/toyredis/internal/resp"

type Dispatcher interface {
	Dispatch(bulkArray resp.BulkArray) []byte
}

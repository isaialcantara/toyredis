package storage

type StringStorage interface {
	Get(key string) ([]byte, error)
	Set(key string, val []byte) error
	Del(key string) (bool, error)
}

package cache

type Cache interface {
	GetUint64(key string) (uint64, error)
	GetByteBuffer(key string) ([]byte, error)

	SetUint64(key string, value uint64)
	SetByteBuffer(key string, value []byte)
}

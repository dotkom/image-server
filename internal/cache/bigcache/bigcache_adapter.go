package bigcache_adapter

import (
	"encoding/binary"
	"time"

	"github.com/allegro/bigcache"
	log "github.com/sirupsen/logrus"
)

type BigCacheAdapter struct {
	bc *bigcache.BigCache
}
type CacheMissError struct{}

func New() (*BigCacheAdapter, error) {
	adapter := &BigCacheAdapter{}
	var err error
	adapter.bc, err = bigcache.NewBigCache(bigcache.DefaultConfig(10 * time.Minute))
	if err != nil {
		return nil, err
	}
	return adapter, nil
}

func (adapter *BigCacheAdapter) GetUint64(key string) (uint64, error) {
	buffer, err := adapter.bc.Get(key)
	if err != nil {
		return 0, &CacheMissError{}
	}

	return binary.BigEndian.Uint64(buffer), nil
}

func (adapter *BigCacheAdapter) GetByteBuffer(key string) ([]byte, error) {
	buffer, err := adapter.bc.Get(key)
	if err != nil {
		return nil, &CacheMissError{}
	}

	return buffer, nil
}

func (adapter *BigCacheAdapter) SetUint64(key string, value uint64) {
	buffer := make([]byte, 8)
	binary.BigEndian.PutUint64(buffer, value)
	if err := adapter.bc.Set(key, buffer); err != nil {
		log.Error(err)
	}
}

func (adapter *BigCacheAdapter) SetByteBuffer(key string, value []byte) {
	if err := adapter.bc.Set(key, value); err != nil {
		log.Error(err)
	}
}

func (err *CacheMissError) Error() string {
	return "cache miss"
}

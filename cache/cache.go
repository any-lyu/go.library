package cache

import (
	"sync"
	"time"

	"github.com/allegro/bigcache"

	"github.com/any-lyu/go.library/logs"
)

var (
	once  sync.Once
	cache *bigcache.BigCache
)

func getCache() *bigcache.BigCache {
	once.Do(func() {
		var err error
		cache, err = bigcache.NewBigCache(bigcache.Config{
			Shards:             16,
			LifeWindow:         24 * 30 * time.Hour,
			MaxEntriesInWindow: 1000 * 10,
			MaxEntrySize:       256,
			//CleanWindow:        61 * time.Minute,
			Verbose: true,
		})
		if err != nil {
			logs.Error("cache init err:", err)
		}
	})
	return cache
}

// Set Cache set
func Set(key string, entry []byte) error {
	return getCache().Set(key, entry)
}

// Get Cache get
func Get(key string) ([]byte, error) {
	return getCache().Get(key)
}

// Delete Cache Delete
func Delete(key string) error {
	return getCache().Delete(key)
}

// Reset Cache Reset
func Reset() error {
	return getCache().Reset()
}

// IsNotFount Cache err is not fount
func IsNotFount(err error) bool {
	return err == bigcache.ErrEntryNotFound
}

package xfreecache

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/coocood/freecache"
)

type (
	Cache = freecache.Cache
)

const (
	defaultSize = 10 * 1024 * 1024
)

var (
	globalCache *Cache
	initOnce    sync.Once
)

type Opt struct {
	Size int
}

func init() {
	globalCache = freecache.NewCache(defaultSize)
}

func Init(opt *Opt) {
	initOnce.Do(func() {
		if opt == nil {
			opt = &Opt{}
		}
		if opt.Size == 0 {
			opt.Size = defaultSize
		}
		globalCache = freecache.NewCache(opt.Size)
	})
}

func getKey(key string) []byte {
	return []byte(key)
}

func Set(key string, value any, expiredTimes ...time.Duration) error {
	expireSeconds := 0
	if len(expiredTimes) > 0 {
		expireSeconds = int(expiredTimes[0].Seconds())
	}
	k := getKey(key)
	v, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return globalCache.Set(k, v, expireSeconds)
}

func Get(key string, value any) (bool, error) {
	k := getKey(key)
	v, err := globalCache.Get(k)
	if err != nil {
		return false, err
	}
	return true, json.Unmarshal(v, &value)
}

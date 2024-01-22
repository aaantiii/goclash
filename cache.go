package clash

import (
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"
	cmap "github.com/orcaman/concurrent-map/v2"
)

type cachedValue struct {
	data  []byte      // data is the cached value
	timer *time.Timer // timer schedules removal of cached value
}

var (
	useCache       = true
	cache          = cmap.New[*cachedValue]()
	fixedCacheTime time.Duration
)

// UseCache sets whether to use cache.
//
// Cache is capable of storing large amounts of data in memory, across different shards. Using it together with UseFixedCacheTime may replace the need for a store like Redis, depending on your needs.
//
// Important: set this before using Client, don't change during runtime.
func UseCache(v bool) {
	useCache = v
}

// UseFixedCacheTime sets a fixed cache time, ignoring Cache-Control headers. Disable by passing 0 as argument.
//
// Important: set this before using Client, don't change during runtime.
func UseFixedCacheTime(d time.Duration) {
	fixedCacheTime = d
}

func readCache(key string) ([]byte, bool) {
	if !useCache {
		return nil, false
	}
	value, ok := cache.Get(key)
	return value.data, ok
}

func writeCache(key string, data []byte, duration time.Duration) {
	if !useCache {
		return
	}
	if value, ok := cache.Get(key); ok {
		value.timer.Stop()
	}
	cache.Set(key, &cachedValue{
		data: data,
		timer: time.AfterFunc(duration, func() {
			cache.Remove(key)
		}),
	})
}

func cacheResponse(url string, res *resty.Response) {
	if fixedCacheTime > 0 {
		writeCache(url, res.Body(), fixedCacheTime)
		return
	}

	seconds, err := strconv.Atoi(res.Header().Get("Cache-Control")[8:])
	if err != nil {
		return
	}
	writeCache(url, res.Body(), time.Duration(seconds)*time.Second)
}

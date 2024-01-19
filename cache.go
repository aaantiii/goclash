package clash

import (
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"
	cmap "github.com/orcaman/concurrent-map/v2"
)

type cachedValue struct {
	data  []byte      // data is the cached value
	timer *time.Timer // timer to remove value from cache
}

var (
	useCache = true
	cache    = cmap.New[*cachedValue]()
)

// UseCache sets whether to use cache. Set this before any request, don't change during runtime.
func UseCache(v bool) {
	useCache = v
}

func readCache(key string) ([]byte, bool) {
	if !useCache {
		return nil, false
	}
	value, ok := cache.Get(key)
	return value.data, ok
}

func writeCache(key string, data []byte, duration int) {
	if !useCache {
		return
	}
	if value, ok := cache.Get(key); ok {
		value.timer.Stop()
	}
	cache.Set(key, &cachedValue{
		data: data,
		timer: time.AfterFunc(time.Duration(duration), func() {
			cache.Remove(key)
		}),
	})
}

func cacheResponse(url string, res *resty.Response) error {
	dur, err := strconv.Atoi(res.Header().Get("Cache-Control"))
	if err != nil {
		return err
	}
	writeCache(url, res.Body(), dur)
	return nil
}

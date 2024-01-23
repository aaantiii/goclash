package goclash

import (
	"strconv"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
	cmap "github.com/orcaman/concurrent-map/v2"
)

// Cache is a simple but performant in-memory cache.
type Cache struct {
	enabled   bool
	store     cmap.ConcurrentMap[string, *cachedValue]
	cacheTime time.Duration
	mu        sync.Mutex
}

type cachedValue struct {
	data  []byte      // data is the cached value
	timer *time.Timer // timer schedules removal of cached value
}

func newCache() *Cache {
	return &Cache{
		store:   cmap.New[*cachedValue](),
		enabled: true,
	}
}

// UseCache sets whether to use cache.
//
// Cache is capable of storing large amounts of data in memory, across different shards. Using it together with SetCacheTime may replace the need for a store like Redis, depending on your needs.
func (h *Client) UseCache(v bool) {
	h.cache.mu.Lock()
	defer h.cache.mu.Unlock()
	h.cache.enabled = v
}

// SetCacheTime sets a fixed cache time, ignoring CacheControl headers. Disable by passing 0 as argument.
func (h *Client) SetCacheTime(d time.Duration) {
	h.cache.mu.Lock()
	defer h.cache.mu.Unlock()
	h.cache.cacheTime = d
}

// Get gets a value from the cache, and a boolean indicating whether the value was found.
func (c *Cache) Get(key string) ([]byte, bool) {
	if !c.enabled {
		return nil, false
	}

	value, ok := c.store.Get(key)
	if !ok {
		return nil, false
	}
	return value.data, ok
}

// Set sets a value in the cache, with a duration after it gets removed.
func (c *Cache) Set(key string, data []byte, duration time.Duration) {
	if !c.enabled {
		return
	}
	if value, ok := c.store.Get(key); ok {
		value.timer.Stop()
	}
	c.store.Set(key, &cachedValue{
		data: data,
		timer: time.AfterFunc(duration, func() {
			c.store.Remove(key)
		}),
	})
}

// CacheResponse caches the response body of a resty.Response, using the Cache-Control header to determine the cache time.
func (c *Cache) CacheResponse(url string, res *resty.Response) {
	if c.cacheTime > 0 {
		c.Set(url, res.Body(), c.cacheTime)
		return
	}
	seconds, err := strconv.Atoi(res.Header().Get("Cache-Control")[8:])
	if err != nil {
		return
	}
	c.Set(url, res.Body(), time.Duration(seconds)*time.Second)
}

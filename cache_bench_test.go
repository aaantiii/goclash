package goclash

import (
	"crypto/rand"
	"strconv"
	"testing"
	"time"

	cmap "github.com/orcaman/concurrent-map/v2"
)

func randomBytes(size int) []byte {
	b := make([]byte, size)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	return b
}

var cacheBenchmarkArgs = []struct {
	name string
	size int
}{
	{"1KB", 1024},
	{"10KB", 1024 * 10},
	{"100KB", 1024 * 100},
}

func BenchmarkCacheReadPtr(b *testing.B) {
	c := cmap.New[*cachedValue]()
	timer := time.NewTimer(time.Minute)
	for _, args := range cacheBenchmarkArgs {
		for i := 0; i < b.N; i++ {
			c.Set(strconv.Itoa(i), &cachedValue{
				data:  randomBytes(args.size),
				timer: timer,
			})
		}
	}

	for _, args := range cacheBenchmarkArgs {
		b.Run(args.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				c.Get(strconv.Itoa(i))
			}
		})
	}
}

func BenchmarkCacheReadVal(b *testing.B) {
	c := cmap.New[cachedValue]()
	timer := time.NewTimer(time.Minute)
	for _, args := range cacheBenchmarkArgs {
		for i := 0; i < b.N; i++ {
			c.Set(strconv.Itoa(i), cachedValue{
				data:  randomBytes(args.size),
				timer: timer,
			})
		}
	}

	for _, args := range cacheBenchmarkArgs {
		b.Run(args.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				c.Get(strconv.Itoa(i))
			}
		})
	}
}

func BenchmarkCacheWritePtr(b *testing.B) {
	c := cmap.New[*cachedValue]()
	timer := time.NewTimer(time.Minute)
	for _, args := range cacheBenchmarkArgs {
		b.Run(args.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				c.Set("key", &cachedValue{
					data:  randomBytes(args.size),
					timer: timer,
				})
			}
		})
	}
}

func BenchmarkCacheWriteVal(b *testing.B) {
	c := cmap.New[cachedValue]()
	timer := time.NewTimer(time.Minute)
	for _, args := range cacheBenchmarkArgs {
		b.Run(args.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				c.Set(strconv.Itoa(i), cachedValue{
					data:  randomBytes(args.size),
					timer: timer,
				})
			}
		})
	}
}

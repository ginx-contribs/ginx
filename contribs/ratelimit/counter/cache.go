package counter

import (
	"fmt"
	"github.com/patrickmn/go-cache"
	"golang.org/x/net/context"
	"sync"
	"time"
)

func Cache() *CacheCounter {
	c := cache.New(cache.NoExpiration, cache.NoExpiration)
	return &CacheCounter{cache: c}
}

// CacheCounter implements Counter by go-cache, and ensuring atomic through sync.Mutex.
type CacheCounter struct {
	mu    sync.Mutex
	cache *cache.Cache
}

func (c *CacheCounter) Count(ctx context.Context, key string, limit int, window time.Duration) (int, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	count, ok := c.cache.Get(key)
	if !ok {
		c.cache.Set(key, 1, window)
	} else if count.(int) < limit {
		c.cache.IncrementInt(key, 1)
	}
	get, _ := c.cache.Get(key)
	fmt.Println(get)
	return get.(int), nil
}

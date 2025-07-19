package pokecache

import (
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

type cache struct {
	CacheEntries map[string]cacheEntry
	mu           sync.RWMutex
}

func NewCache(interval time.Duration) *cache {
	newCache := &cache{
		CacheEntries: map[string]cacheEntry{},
	}
	go newCache.reapLoop(interval)
	return newCache
}

func (c *cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.CacheEntries[key] = cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
}

func (c *cache) Get(key string) ([]byte, bool) {
	c.mu.RLock()
	entry, ok := c.CacheEntries[key]
	c.mu.RUnlock()
	if ok {
		return entry.val, ok
	}
	return nil, ok
}

func (c *cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		<-ticker.C
		c.mu.Lock()
		for key, entry := range c.CacheEntries {
			if time.Since(entry.createdAt) > interval {
				delete(c.CacheEntries, key)
			}
		}
		c.mu.Unlock()
	}
}

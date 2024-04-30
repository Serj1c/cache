package main

import (
	"sync"
	"time"
)

type cacheItem struct {
	value interface{}
	exp   time.Time
}

type CacheWithLazyTtl struct {
	data map[string]cacheItem
	mu   sync.RWMutex
}

func NewCacheWithLazyTtl() *CacheWithLazyTtl {
	return &CacheWithLazyTtl{
		data: map[string]cacheItem{},
	}
}

func (c *CacheWithLazyTtl) Get(key string) (value interface{}, ok bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, ok := c.data[key]
	if !ok {
		return nil, false
	}

	if item.exp.Before(time.Now()) {
		delete(c.data, key)
		return nil, false
	}

	return item.value, true
}

func (c *CacheWithLazyTtl) Set(key string, val interface{}, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data[key] = cacheItem{
		value: val,
		exp:   time.Now().Add(ttl),
	}
}

func (c *CacheWithLazyTtl) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.data, key)
}

func (c *CacheWithLazyTtl) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data = map[string]cacheItem{}
}

package main

import (
	"sync"
	"time"
)

type cacheItem struct {
	value interface{}
	exp   int64
}

type CacheWithTtl struct {
	data              map[string]cacheItem
	defaultExpiration time.Duration
	gcInterval        time.Duration
	mu                sync.RWMutex
}

func NewCacheWithTtl(defaultExpiration, gcInterval time.Duration) *CacheWithTtl {
	c := &CacheWithTtl{
		data:              map[string]cacheItem{},
		defaultExpiration: defaultExpiration,
		gcInterval:        gcInterval,
	}

	if gcInterval > 0 {
		go c.gcRun()
	}

	return c
}

func (c *CacheWithTtl) gcRun() {
	if keys := c.getExpiredKeys(); len(keys) > 0 {
		c.deleteExpiredKeys(keys)
	}
}

func (c *CacheWithTtl) Get(key string) (value interface{}, ok bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, ok := c.data[key]
	if !ok {
		return nil, false
	}

	if item.exp > 0 {
		if time.Now().UnixNano() > item.exp {
			return nil, false
		}
	}

	return item.value, true
}

func (c *CacheWithTtl) Set(key string, val interface{}, ttl time.Duration) {
	var exp int64
	if ttl == 0 {
		ttl = c.defaultExpiration
	}

	if ttl > 0 {
		exp = time.Now().Add(ttl).UnixNano()
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.data[key] = cacheItem{
		value: val,
		exp:   exp,
	}
}

func (c *CacheWithTtl) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.data, key)
}

func (c *CacheWithTtl) getExpiredKeys() []string {
	keys := make([]string, 0, 0)
	c.mu.RLock()
	defer c.mu.RUnlock()
	for k, v := range c.data {
		if time.Now().UnixNano() > v.exp && v.exp > 0 {
			keys = append(keys, k)
		}
	}
	return keys
}

func (c *CacheWithTtl) deleteExpiredKeys(keys []string) []string {
	c.mu.Lock()
	defer c.mu.Unlock()
	for _, k := range keys {
		delete(c.data, k)
	}
	return keys
}

func (c *CacheWithTtl) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data = map[string]cacheItem{}
}

func (c *CacheWithTtl) clearItems() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data = map[string]cacheItem{}
}

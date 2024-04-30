package main

import "sync"

type SimpleCache struct {
	data map[string]interface{}
	mu   sync.RWMutex
}

func NewCache() *SimpleCache {
	return &SimpleCache{
		data: map[string]interface{}{},
	}
}

func (c *SimpleCache) Get(key string) (value interface{}, ok bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	val, ok := c.data[key]
	return val, ok
}

func (c *SimpleCache) Set(key string, val interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data[key] = val
}

func (c *SimpleCache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.data, key)
}

func (c *SimpleCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data = map[string]interface{}{}
}

package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

type cache struct {
	mu    sync.RWMutex
	value int
}

func newCache() *cache {
	c := &cache{
		value: longBoringCalc(),
	}
	go c.update()
	return c
}

func (c *cache) update() {
	for range time.Tick(time.Second * 2) {
		val := longBoringCalc()
		c.mu.Lock()
		c.value = val
		c.mu.Unlock()
	}
}

func longBoringCalc() int {
	time.Sleep(time.Second)
	return rand.Intn(1000)
}

func (c *cache) get() int {
	c.mu.RLock()
	val := c.value
	c.mu.RUnlock()
	return val
}

func main() {
	c := newCache()
	http.HandleFunc("/calc", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "{\"result\":%d}", c.get())
	})

	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}

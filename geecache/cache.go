package geecache

import (
	"sync"
)

type Cache struct {
	mu         sync.Mutex
	inner      cacheType
	cacheBytes int
}

func NewCache(cacheBytes int, onEvicted func(key string, value interface{}), t T) *Cache {
	ty := t.String()
	var inner cacheType
	if ty == "fifo" {
		inner = NewFIFO(cacheBytes, onEvicted)
	} else if ty == "lfu" {
		inner = NewLFU(cacheBytes, onEvicted)
	} else {
		inner = NewLRU(cacheBytes, onEvicted)
	}
	return &Cache{
		inner:      inner,
		cacheBytes: cacheBytes,
	}
}

func (c *Cache) add(key string, value ByteView) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.inner.Put(key, value)
}

func (c *Cache) get(key string) (value ByteView, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if v, ok:= c.inner.Get(key); ok {
		return v.(ByteView), ok
	}

	return
}

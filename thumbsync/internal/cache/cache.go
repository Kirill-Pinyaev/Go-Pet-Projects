package cache

import "sync"

// Cache хранит миниатюры в памяти.

type Cache struct {
	mu  sync.RWMutex
	mem map[string][]byte
}

func New() *Cache {
	return &Cache{mem: make(map[string][]byte)}
}

func (c *Cache) Load(url string) ([]byte, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	v, ok := c.mem[url]
	return v, ok
}

func (c *Cache) Store(url string, data []byte) {
	c.mu.Lock()
	c.mem[url] = data
	c.mu.Unlock()
}

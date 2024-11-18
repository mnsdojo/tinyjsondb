package internal

import (
	"sync"
	"time"
)

// cache structure to hold data and its expiration time
type Cache struct {
	data       map[string]cachedItem
	mu         sync.RWMutex
	expiration time.Duration
}

// cachedItem holds data and its expiration time.
type cachedItem struct {
	value      interface{}
	expiration time.Time
}

func NewCache(expiration time.Duration) *Cache {
	return &Cache{
		data:       make(map[string]cachedItem),
		expiration: expiration,
	}
}

// Get retrieves a cached item by key. If the item is expired or doesn't exist, it returns nil.
func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, found := c.data[key]
	if !found || (c.expiration > 0 && time.Now().After(item.expiration)) {
		return nil, false
	}
	return item.value, true
}

func (c *Cache) Set(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	var expirationTime time.Time

	if c.expiration > 0 {
		expirationTime = time.Now().Add(c.expiration)
	}
	c.data[key] = cachedItem{
		value:      key,
		expiration: expirationTime,
	}
}

func (c *Cache) Invalidate(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.data, key)
}

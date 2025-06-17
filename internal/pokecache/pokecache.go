package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	items		map[string]cacheEntry
	mutex		*sync.RWMutex
}

type cacheEntry struct {
	createdAt	time.Time
	val			[]byte
}


func NewCache(interval time.Duration) *Cache {
	cache := Cache {
		map[string]cacheEntry{},
		&sync.RWMutex{},
	}

	go cache.reapLoop(interval)
	return &cache
}

func (c *Cache) Add(key string, val []byte) {
	c.mutex.Lock()
	c.items[key] = cacheEntry {
		time.Now(),
		val,
	};
	c.mutex.Unlock()
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mutex.RLock()
	item, ok := c.items[key]
	c.mutex.RUnlock()

	if !ok {
		return nil, false
	}

	return item.val, true
}

func (c *Cache) reapLoop(dur time.Duration) {
	ticker := time.NewTicker(dur)

	for {
		<- ticker.C
		for key,item := range c.items {
			expiry := item.createdAt.Add(dur)
			if expiry.Compare(time.Now()) == -1 {
				c.mutex.Lock()
				delete(c.items, key) 
				c.mutex.Unlock()
			}
		}
	}
}

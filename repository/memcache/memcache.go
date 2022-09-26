package memcache

import (
	"GoVoteApi/repository"
	"sync"
)

type cache struct {
    cacheMap map[string]any
    sync.Mutex
}

func (c *cache) Set(key string, val any)  {
    c.Lock()
    defer c.Unlock()
    c.cacheMap[key] = val
}

// Get return nil if key not found
func (c *cache) Get(key string) any {
    c.Lock()
    defer c.Unlock()
    if v, ok := c.cacheMap[key]; ok {
        return v
    }
    return nil 
}

func New() repository.Cache {
    return &cache{
        cacheMap: make(map[string]any),
        Mutex: sync.Mutex{},
    }
}

package cache

import (
	"distributedCache/pkg/cache/policy"
	"sync"
	"time"
)

type PolicyType int

const (
	LRU PolicyType = iota
	LFU
)

type CacheItem[V any] struct {
	Value      V
	Expiration int64
}

type Cache[K comparable, V any] struct {
	capacity int
	policy   policy.EvictionPolicy[K, V]
	data     map[K]*CacheItem[V]
	mu       sync.RWMutex
}

func NewCache[K comparable, V any](capacity int, cleanupInterval time.Duration, policyType PolicyType) *Cache[K, V] {
	var cacheEvictPolicy policy.EvictionPolicy[K, V]
	// Select the appropriate policy
	switch policyType {
	case LRU:
		cacheEvictPolicy = policy.NewLRU[K, V](capacity)
	case LFU:
		cacheEvictPolicy = policy.NewLFU[K, V](capacity)
	default:
		panic("unknown cache policy")
	}

	return &Cache[K, V]{
		capacity: capacity,
		policy:   cacheEvictPolicy,
		data:     make(map[K]*CacheItem[V]),
	}
}

func (c *Cache[K, V]) Set(key K, value V, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	expiration := time.Now().Add(ttl).UnixNano()
	c.data[key] = &CacheItem[V]{
		Value:      value,
		Expiration: expiration,
	}

	c.policy.Add(key, value)

	if len(c.data) > c.capacity {
		if key, _, ok := c.policy.Evict(); ok {
			delete(c.data, key)
		}
	}
}

func (c *Cache[K, V]) Get(key K) (V, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	item, exist := c.data[key]
	if !exist || item.Expiration < time.Now().UnixNano() {
		if exist {
			delete(c.data, key)
			c.policy.Remove(key)
		}
		var zero V
		return zero, false
	}

	c.policy.RecordAccess(key, item.Value)
	return item.Value, true
}

func (c *Cache[K, V]) Delete(key K) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	_, exist := c.data[key]
	if exist {
		delete(c.data, key)
		c.policy.Remove(key)
		return true
	}
	return false
}

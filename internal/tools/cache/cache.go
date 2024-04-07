package cache

import (
	"context"
	"sync"
	"time"
)

type CacheSettings struct {
	TTL             time.Duration
	CleanupInterval time.Duration
}

type Cache[K comparable, V any] struct {
	Data map[K]V

	Exp map[K]time.Time

	Expiration time.Duration

	Mu sync.Mutex
}

func New[K comparable, V any](ctx context.Context, settings *CacheSettings) *Cache[K, V] {
	cache := &Cache[K, V]{
		Data: make(map[K]V),
		Mu:   sync.Mutex{},
		Exp:  make(map[K]time.Time),
	}

	go cache.Cleanup(ctx, settings.CleanupInterval)

	return cache
}

func (c *Cache[K, V]) Cleanup(ctx context.Context, interval time.Duration) {
	t := time.NewTicker(interval)
	defer t.Stop()

	for {
		select {
		case <-t.C:
			c.deleteExpired()

		case <-ctx.Done():
			return
		}
	}
}

func (c *Cache[K, V]) Get(key K) (value V, ok bool) {
	c.Mu.Lock()
	defer c.Mu.Unlock()

	val, ok := c.Data[key]
	if !ok {
		return
	}

	if expired(c.Exp[key]) {
		return value, false
	}

	return val, true
}

func (c *Cache[K, V]) Set(key K, val V) {
	c.Mu.Lock()

	c.Data[key] = val
	c.Exp[key] = time.Now().Add(c.Expiration)

	c.Mu.Unlock()
}

func (c *Cache[K, V]) Delete(key K) {
	c.Mu.Lock()

	delete(c.Data, key)
	delete(c.Exp, key)

	c.Mu.Unlock()
}

func (c *Cache[K, V]) deleteExpired() {
	c.Mu.Lock()

	for key, exp := range c.Exp {
		if expired(exp) {
			delete(c.Data, key)
			delete(c.Exp, key)
		}
	}

	c.Mu.Unlock()
}

package tools

import (
	"context"
	"reflect"
	"sync"
	"time"

	"github.com/v1shn3vsk7/mongorm/internal/tools/cache"
)

type ExternalTools struct {
	Cache *cache.Cache[any, any]
}

func (t *ExternalTools) CreateCache(ctx context.Context, key any, TTL, cleanup time.Duration) {
	tp := reflect.TypeOf(key)
	if !tp.Comparable() {
		panic("received not comparable key")
	}

	t.Cache = &cache.Cache[any, any]{
		Data:       make(map[any]any),
		Exp:        make(map[any]time.Time),
		Expiration: TTL,
		Mu:         sync.Mutex{},
	}

	go t.Cache.Cleanup(ctx, cleanup)
}

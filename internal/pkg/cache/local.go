package cache

import (
	"context"
	"encoding/json"
	"github.com/allegro/bigcache/v3"
	"time"
)

const (
	// DefaultExpirationTime is the default expiration time for cache entries
	DefaultExpirationTime = 600
	CleanWindowMinutes    = 5
)

type LocalCache interface {
	Get(key string, value interface{}) (bool, error)
	Set(key string, value interface{}) error
	Delete(key string) error
}

type localCache struct {
	cache *bigcache.BigCache
}

func NewLocalCache(ctx context.Context, expirationTime uint64) (LocalCache, error) {
	if expirationTime == 0 {
		expirationTime = DefaultExpirationTime
	}
	cacheCfg := bigcache.DefaultConfig(time.Duration(expirationTime) * time.Second)
	cacheCfg.CleanWindow = time.Duration(CleanWindowMinutes) * time.Minute
	cache, err := bigcache.New(ctx, cacheCfg)
	if err != nil {
		return nil, err
	}

	return &localCache{cache: cache}, nil
}

func (c *localCache) Get(key string, value interface{}) (bool, error) {
	val, err := c.cache.Get(key)
	if err == bigcache.ErrEntryNotFound {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	if err = json.Unmarshal(val, value); err != nil {
		return false, err
	}

	return true, nil
}

func (c *localCache) Set(key string, value interface{}) error {
	val, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return c.cache.Set(key, val)
}

func (c *localCache) Delete(key string) error {
	return c.cache.Delete(key)
}

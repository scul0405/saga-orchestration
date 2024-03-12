package cache

import (
	"context"
	"encoding/json"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	"github.com/redis/go-redis/v9"
	"time"
)

type RedisCache interface {
	Get(ctx context.Context, key string, value interface{}) (bool, error)
	Set(ctx context.Context, key string, value interface{}) error
	Delete(ctx context.Context, key string) error
	CFReserve(ctx context.Context, key string, capacity, bucketSize, maxIterations int64) error
	CFExist(ctx context.Context, key string, value interface{}) (bool, error)
	CFAdd(ctx context.Context, key string, value interface{}) error
	CFDel(ctx context.Context, key string, value interface{}) error
	GetMutex(mutexName string) *redsync.Mutex
	ExecIncrbyXPipeline(ctx context.Context, payloads *[]RedisIncrbyXPayload) error
}

type RedisIncrbyXPayload struct {
	Key   string
	Value int64
}

type redisCache struct {
	client         redis.UniversalClient
	rs             *redsync.Redsync
	expirationTime time.Duration
}

func NewRedisCache(client redis.UniversalClient, expirationTime time.Duration) RedisCache {
	pool := goredis.NewPool(client)
	rs := redsync.New(pool)

	return &redisCache{
		client:         client,
		expirationTime: expirationTime,
		rs:             rs,
	}
}

func (c *redisCache) Get(ctx context.Context, key string, value interface{}) (bool, error) {
	val, err := c.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	if err = json.Unmarshal([]byte(val), value); err != nil {
		return false, err
	}

	return true, nil
}

func (c *redisCache) Set(ctx context.Context, key string, value interface{}) error {
	val, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return c.client.Set(ctx, key, val, c.expirationTime).Err()
}

func (c *redisCache) Delete(ctx context.Context, key string) error {
	return c.client.Del(ctx, key).Err()
}

func (c *redisCache) CFReserve(ctx context.Context, key string, capacity, bucketSize, maxIterations int64) error {
	return c.client.CFReserveWithArgs(ctx, key, &redis.CFReserveOptions{
		Capacity:      capacity,
		BucketSize:    bucketSize,
		MaxIterations: maxIterations,
	}).Err()
}

func (c *redisCache) CFExist(ctx context.Context, key string, value interface{}) (bool, error) {
	val, err := c.client.CFExists(ctx, key, value).Result()
	if err != nil {
		return false, err
	}
	return val, nil
}

func (c *redisCache) CFAdd(ctx context.Context, key string, value interface{}) error {
	return c.client.CFAdd(ctx, key, value).Err()
}

func (c *redisCache) CFDel(ctx context.Context, key string, value interface{}) error {
	return c.client.CFDel(ctx, key, value).Err()
}

func (c *redisCache) GetMutex(mutexName string) *redsync.Mutex {
	return c.rs.NewMutex(mutexName)
}

// Lua script to increment a key by a value if it exists
var incrByX = redis.NewScript(`
local exists = redis.call("EXISTS", KEYS[1])
if exists == 1 then
	return redis.call("INCRBY", KEYS[1], ARGV[1])
end
`)

func (c *redisCache) ExecIncrbyXPipeline(ctx context.Context, payloads *[]RedisIncrbyXPayload) error {
	pipe := c.client.Pipeline()
	executedCmds := make([]*redis.Cmd, len(*payloads))
	for i, payload := range *payloads {
		executedCmds[i] = incrByX.Run(ctx, c.client, []string{payload.Key}, payload.Value)
	}
	_, err := pipe.Exec(ctx)
	if err != nil {
		return err
	}

	for _, cmd := range executedCmds {
		if err = cmd.Err(); err != nil {
			return err
		}
	}

	return nil
}

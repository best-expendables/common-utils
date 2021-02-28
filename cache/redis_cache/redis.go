package redis_cache

import (
	"context"
	"time"

	"github.com/best-expendables/common-utils/cache"

	redisCache "github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
)

type Redis struct {
	Client *redis.Client
	Cache  *redisCache.Cache
	Prefix string
	Ttl    time.Duration
}

func NewRedis(c *redis.Client, prefix string, ttl time.Duration) *Redis {
	return &Redis{
		Client: c,
		Cache:  NewRedisCache(c),
		Prefix: prefix,
		Ttl:    ttl,
	}
}

func NewRedisCacheCreateFunc(prefix string, ttl time.Duration) func(c *redis.Client) *Redis {
	return func(c *redis.Client) *Redis {
		return &Redis{
			Client: c,
			Cache:  NewRedisCache(c),
			Prefix: prefix,
			Ttl:    ttl,
		}
	}
}

func (r *Redis) Get(ctx context.Context, key string, obj interface{}) error {
	err := r.Cache.Get(ctx, r.cacheKey(key), obj)
	if err == redisCache.ErrCacheMiss {
		return cache.Nil
	}
	return err
}

func (r *Redis) MGet(ctx context.Context, keys ...string) ([]interface{}, error) {
	result, err := r.Client.MGet(ctx, r.cacheKeys(keys...)...).Result()
	for i := range result {
		if result[i] == nil {
			return nil, cache.Nil
		}
	}
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *Redis) Set(ctx context.Context, key string, obj interface{}) error {
	return r.Cache.Set(&redisCache.Item{
		Ctx:   ctx,
		Key:   r.cacheKey(key),
		Value: obj,
		TTL:   r.Ttl,
	})
}

func (r *Redis) HSet(ctx context.Context, key string, field string, obj interface{}) error {
	_, err := r.Client.HSet(ctx, r.cacheKey(key), field, obj).Result()

	return err
}

func (r *Redis) HGet(ctx context.Context, key string, field string) (string, error) {
	hGet := r.Client.HGet(ctx, r.cacheKey(key), field)
	if err := hGet.Err(); err != nil {
		return "", err
	}

	return hGet.Val(), nil
}

func (r *Redis) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	hGetAll := r.Client.HGetAll(ctx, r.cacheKey(key))
	if err := hGetAll.Err(); err != nil {
		return nil, err
	}

	return hGetAll.Val(), nil
}

func (r *Redis) MSet(ctx context.Context, obj ...interface{}) error {
	return r.Client.MSet(ctx, obj...).Err()
}

func (r *Redis) Delete(ctx context.Context, key string) error {
	err := r.Cache.Delete(ctx, r.cacheKey(key))
	if err == redisCache.ErrCacheMiss {
		return cache.Nil
	}
	return err
}

func (r *Redis) ScanD(ctx context.Context, match string) error {
	scan := r.Client.Scan(ctx, 0, r.cacheKey(match), 0).Iterator()
	for scan.Next(ctx) {
		err := r.Client.Del(ctx, scan.Val()).Err()
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *Redis) HExpire(ctx context.Context, key string) error {
	_, err := r.Client.Expire(ctx, r.cacheKey(key), r.Ttl).Result()
	return err
}

func (r *Redis) cacheKey(key string) string {
	return r.Prefix + "/" + key
}

func (r *Redis) cacheKeys(keys ...string) []string {
	for i := range keys {
		keys[i] = r.Prefix + "/" + keys[i]
	}
	return keys
}

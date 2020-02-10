package redis_cache

import (
	"bitbucket.org/snapmartinc/common-utils/cache"
	"time"

	redisCache "github.com/go-redis/cache"
	"github.com/go-redis/redis"
)

type Redis struct {
	client *redis.Client
	codec  *redisCache.Codec
	prefix string
	ttl    time.Duration
}

func NewRedis(c *redis.Client, prefix string, ttl time.Duration) *Redis {
	return &Redis{
		client: c,
		codec:  NewRedisCodec(c),
		prefix: prefix,
		ttl:    ttl,
	}
}

func NewRedisCacheCreateFunc(prefix string, ttl time.Duration) func(c *redis.Client) *Redis {
	return func(c *redis.Client) *Redis {
		return &Redis{
			client: c,
			codec:  NewRedisCodec(c),
			prefix: prefix,
			ttl:    ttl,
		}
	}
}

func (r *Redis) Get(key string, obj interface{}) error {
	err := r.codec.Get(r.cacheKey(key), obj)
	if err == redisCache.ErrCacheMiss {
		return cache.Nil
	}
	return err
}

func (r *Redis) MGet(keys ...string) ([]interface{}, error) {
	result, err := r.client.MGet(r.cacheKeys(keys...)...).Result()
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

func (r *Redis) Set(key string, obj interface{}) error {
	return r.codec.Set(&redisCache.Item{
		Key:        r.cacheKey(key),
		Object:     obj,
		Expiration: r.ttl,
	})
}

func (r *Redis) MSet(obj ...interface{}) error {
	return r.client.MSet(obj...).Err()
}

func (r *Redis) Delete(key string) error {
	err := r.codec.Delete(r.cacheKey(key))
	if err == redisCache.ErrCacheMiss {
		return cache.Nil
	}
	return err
}

func (r *Redis) cacheKey(key string) string {
	return r.prefix + "/" + key
}

func (r *Redis) cacheKeys(keys ...string) []string {
	for i := range keys {
		keys[i] = r.prefix + "/" + keys[i]
	}
	return keys
}

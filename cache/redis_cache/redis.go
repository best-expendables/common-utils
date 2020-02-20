package redis_cache

import (
	"bitbucket.org/snapmartinc/common-utils/cache"
	"time"

	redisCache "github.com/go-redis/cache"
	"github.com/go-redis/redis"
)

type Redis struct {
	Client *redis.Client
	Codec  *redisCache.Codec
	Prefix string
	Ttl    time.Duration
}

func NewRedis(c *redis.Client, prefix string, ttl time.Duration) *Redis {
	return &Redis{
		Client: c,
		Codec:  NewRedisCodec(c),
		Prefix: prefix,
		Ttl:    ttl,
	}
}

func NewRedisCacheCreateFunc(prefix string, ttl time.Duration) func(c *redis.Client) *Redis {
	return func(c *redis.Client) *Redis {
		return &Redis{
			Client: c,
			Codec:  NewRedisCodec(c),
			Prefix: prefix,
			Ttl:    ttl,
		}
	}
}

func (r *Redis) Get(key string, obj interface{}) error {
	err := r.Codec.Get(r.cacheKey(key), obj)
	if err == redisCache.ErrCacheMiss {
		return cache.Nil
	}
	return err
}

func (r *Redis) MGet(keys ...string) ([]interface{}, error) {
	result, err := r.Client.MGet(r.cacheKeys(keys...)...).Result()
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
	return r.Codec.Set(&redisCache.Item{
		Key:        r.cacheKey(key),
		Object:     obj,
		Expiration: r.Ttl,
	})
}

func (r *Redis) HSet(key string, field string, obj interface{}) error {
	_, err := r.Client.HSet(r.cacheKey(key), field, obj).Result()

	return err
}

func (r *Redis) HGet(key string, field string) (string, error) {
	hGet := r.Client.HGet(r.cacheKey(key), field)
	if err := hGet.Err(); err != nil {
		return "", err
	}

	return hGet.Val(), nil
}

func (r *Redis) HGetAll(key string) (map[string]string, error) {
	hGetAll := r.Client.HGetAll(r.cacheKey(key))
	if err := hGetAll.Err(); err != nil {
		return nil, err
	}

	return hGetAll.Val(), nil
}

func (r *Redis) MSet(obj ...interface{}) error {
	return r.Client.MSet(obj...).Err()
}

func (r *Redis) Delete(key string) error {
	err := r.Codec.Delete(r.cacheKey(key))
	if err == redisCache.ErrCacheMiss {
		return cache.Nil
	}
	return err
}

func (r *Redis) ScanD(match string) error {
	scan := r.Client.Scan(0, r.cacheKey(match), 0).Iterator()
	for scan.Next() {
		err := r.Client.Del(scan.Val()).Err()
		if err != nil {
			return err
		}
	}

	return nil
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

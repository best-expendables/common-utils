package redis_cache

import (
	redisCache "github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
)

func NewRedisCache(redis *redis.Client) *redisCache.Cache {
	return redisCache.New(&redisCache.Options{
		Redis:         redis,
		StatsEnabled:  false,
		LocalCache:    nil,
		LocalCacheTTL: 0,
	})
}

package redis_cache

import (
	"encoding/json"

	redisCache "github.com/go-redis/cache/v7"
	"github.com/go-redis/redis/v7"
)

func NewRedisCodec(redis *redis.Client) *redisCache.Codec {
	return &redisCache.Codec{
		Redis:     redis,
		Marshal:   json.Marshal,
		Unmarshal: json.Unmarshal,
	}
}

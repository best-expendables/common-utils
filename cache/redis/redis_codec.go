package redis

import (
	"encoding/json"

	redisCache "github.com/go-redis/cache"
	"github.com/go-redis/redis"
)

func NewRedisCodec(redis *redis.Client) *redisCache.Codec {
	return &redisCache.Codec{
		Redis:     redis,
		Marshal:   json.Marshal,
		Unmarshal: json.Unmarshal,
	}
}

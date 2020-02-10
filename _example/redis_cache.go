package _example

import (
	"bitbucket.org/snapmartinc/common-utils/cache/redis_cache"
	"fmt"
	"github.com/go-redis/redis"
	"time"
)

func main() {
	options := redis.Options{
		Addr:               "redis:6379",
		MaxRetries:         10,
		DialTimeout:        time.Minute,
		ReadTimeout:        time.Minute,
		WriteTimeout:       time.Minute,
		PoolSize:           1000,
		PoolTimeout:        time.Minute,
		IdleTimeout:        time.Minute,
		IdleCheckFrequency: time.Second * 10,
	}

	redisTTL := 60 * 60 * time.Second
	redisConn := redis.NewClient(&options)
	redisClient := redis_cache.NewRedis(redisConn, "redis_prefix", redisTTL)

	var data interface{}

	// Redis set value
	if err := redisClient.Set("redis_key", "redis_value"); err != nil {
		fmt.Println(err)
	}

	// Redis MSet value
	mSetObj := make(map[string]interface{})
	mSetObj["redis_key_001"] = "redis_value_001"
	mSetObj["redis_key_002"] = "redis_value_002"
	if err := redisClient.MSet(mSetObj); err != nil {
		fmt.Println(err)
	}

	// Redis get value
	if err := redisClient.Get("redis_key", &data); err != nil {
		fmt.Println(err)
	}

	// Redis MGet value
	if _, err := redisClient.MGet("redis_key_001", "redis_key_002"); err != nil {
		fmt.Println(err)
	}

	// Redis delete value
	if err := redisClient.Delete("redis_key"); err != nil {
		fmt.Println(err)
	}
}


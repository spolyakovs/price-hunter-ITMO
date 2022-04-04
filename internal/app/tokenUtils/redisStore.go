package tokenUtils

import (
	"os"

	"github.com/go-redis/redis"
)

var redisStore *redis.Client

func SetupRedis() error {
	dsn := os.Getenv("REDIS_DSN")
	if len(dsn) == 0 {
		dsn = "localhost:6379"
	}
	redisStore = redis.NewClient(&redis.Options{
		Addr: dsn,
	})
	_, err := redisStore.Ping().Result()
	return err
}

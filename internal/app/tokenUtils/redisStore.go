package tokenUtils

import (
	"github.com/go-redis/redis"
	"os"
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
	println("+")
	_, err := redisStore.Ping().Result()
	println("++")
	return err
}

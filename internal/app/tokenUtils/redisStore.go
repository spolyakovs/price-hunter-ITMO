package tokenUtils

import (
	"github.com/go-redis/redis"
)

var redisStore *redis.Client

func SetupRedis(redisAddr string) error {
	redisStore = redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})
	_, err := redisStore.Ping().Result()
	return err
}

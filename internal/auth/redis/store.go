package redis

import (
	"github.com/go-redis/redis/v7"
)

// Redis - redis client
type Redis struct {
	Client *redis.Client
	dsn    string
}

// NewClient - ...
func NewClient(dsn string) (*Redis, error) {
	redisClient := &Redis{}
	redisClient.dsn = dsn
	redisClient.Client = redis.NewClient(&redis.Options{
		Addr: dsn,
	})

	_, err := redisClient.Client.Ping().Result()
	if err != nil {
		return nil, err
	}

	return redisClient, nil
}

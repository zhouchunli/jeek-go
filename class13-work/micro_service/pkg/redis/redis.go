package redis

import (
	"context"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

type Conf struct {
	Address     string
	MaxPoolSize int
	Password    string
	DB          int
}

func New(ctx context.Context, conf *Conf) *redis.Client {

	redisOptions := &redis.Options{
		Addr:        conf.Address,
		Password:    conf.Password,
		PoolSize:    conf.MaxPoolSize,
		DB:          conf.DB,
		IdleTimeout: 30 * time.Second,
	}

	redisCache := redis.NewClient(redisOptions)

	pong, err := redisCache.Ping(ctx).Result()
	if err != nil {
		panic(err)
	}

	log.Println("Redis is connected!!!", pong)
	return redisCache
}

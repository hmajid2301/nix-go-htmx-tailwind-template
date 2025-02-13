package redis

import (
	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
)

type Client struct {
	Redis       *redis.Client
	Subscribers map[string]*redis.PubSub
}

func NewRedisClient(address string, retries int) (Client, error) {
	r := redis.NewClient(&redis.Options{
		Addr:       address,
		Password:   "",
		DB:         0,
		MaxRetries: retries,
	})

	if err := redisotel.InstrumentTracing(r); err != nil {
		return Client{}, err
	}

	return Client{
		Redis:       r,
		Subscribers: map[string]*redis.PubSub{},
	}, nil
}

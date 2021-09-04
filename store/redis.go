package store

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

type Redis struct {
	client *redis.Client
}

type ConfigRedis struct {
	Host     string
	Port     int
	Password string
	DB       int
}

func NewRedis(config *ConfigRedis) *Redis {

	log.Println(fmt.Sprintf("%s:%d", config.Host, config.Port))
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Host, config.Port),
		Password: config.Password,
		DB:       config.DB,
	})

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		log.Fatal(err)
	}

	return &Redis{
		client: rdb,
	}
}

func (red *Redis) Get(key string) (string, bool) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)

	defer cancel()

	val, err := red.client.Get(ctx, key).Result()

	if err != nil {
		return "", false
	}

	return val, true
}

func (red *Redis) Set(key string, value string, ttl time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)

	defer cancel()

	return red.client.Set(ctx, key, value, ttl).Err()
}

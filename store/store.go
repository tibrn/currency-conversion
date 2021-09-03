package store

import (
	"currency-conversion/config"
	"time"
)

type Store interface {
	Get(string) (string, bool)
	Set(string, string, time.Duration)
}

var store Store

func init() {

	cfg := config.Get()

	store = NewRedis(&ConfigRedis{
		Host:     cfg.RedisHost,
		Port:     cfg.RedisPort,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})

}

func Get() Store {
	return store
}

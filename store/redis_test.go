package store

import (
	"currency-conversion/config"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func Test_newRedis(t *testing.T) {
	req := require.New(t)

	redis := NewRedis(&ConfigRedis{})

	req.NotNil(redis)
}

func TestRedis_Get_Set(t *testing.T) {

	cfg := config.Get()

	req := require.New(t)

	redis := NewRedis(&ConfigRedis{
		Host:     cfg.RedisHost,
		Port:     cfg.RedisPort,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})

	err := redis.Set("test", "test", time.Millisecond*100)
	req.NoError(err)

	val, isVal := redis.Get("test")

	req.Equal("test", val)
	req.Equal(true, isVal)

	val2, isVal2 := redis.Get("test2")
	req.Equal("", val2)
	req.Equal(false, isVal2)

	time.Sleep(time.Millisecond * 100)

	val2, isVal2 = redis.Get("test")
	req.Equal("", val2)
	req.Equal(false, isVal2)

}

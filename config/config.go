package config

import (
	"log"
	"strconv"
	"time"

	"github.com/gobuffalo/envy"
)

var conf *config

type config struct {
	Port int
	Host string

	FixerAccessKey string

	RedisDB       int
	RedisPort     int
	RedisHost     string
	RedisPassword string

	ExpirationProject time.Duration

	IsDev bool
}

func mustGet(param string) string {
	val, err := envy.MustGet(param)

	if err != nil {
		log.Fatal(err)
	}
	return val
}

func mustGetInt(param string) int {

	nr, err := strconv.Atoi(mustGet(param))

	if err != nil {
		log.Fatal(err)
	}

	return nr

}

func init() {

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	conf = &config{
		Port: mustGetInt("PORT"),
		Host: mustGet("HOST"),

		FixerAccessKey: mustGet("FIXER_ACCESS_KEY"),

		RedisHost:     mustGet("REDIS_HOST"),
		RedisPort:     mustGetInt("REDIS_PORT"),
		RedisPassword: mustGet("REDIS_PASSWORD"),
		RedisDB:       mustGetInt("REDIS_DB"),

		ExpirationProject: time.Hour * 24 * 7,
	}

}

func Get() config {
	return *conf
}

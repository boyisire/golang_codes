package base

import (
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/garyburd/redigo/redis"
	_ "github.com/joho/godotenv/autoload"
)

var RedisMaxActive, _ = strconv.Atoi(os.Getenv("RedisMaxActive"))
var _redis *Redis

type Redis struct {
	Default *redis.Pool
	Live    *redis.Pool
}

/**
 * redis_type [default默认 live直播间]
 */
func NewRedis() *Redis {
	if _redis == nil {
		_redis = &Redis{}
	}
	return _redis
}

func (me *Redis) Open(redis_type string) error {
	switch redis_type {
	case "default":
		me.Default = getRedis(os.Getenv("REDIS_URL"))
	case "live":
		me.Live = getRedis(os.Getenv("LIVE_REDIS_URL"))
	default:
		return errors.New("Redis Type Error!")
	}

	return nil
}

func (me *Redis) Close() {
	if me.Default != nil {
		me.Default.Close()
	}

	if me.Live != nil {
		me.Live.Close()
	}
}

func getRedis(url string) *redis.Pool {
	if RedisMaxActive == 0 {
		RedisMaxActive = 40
	}
	return &redis.Pool{
		MaxIdle:     3,
		MaxActive:   RedisMaxActive,
		IdleTimeout: 240 * time.Second,
		Wait:        true,
		Dial: func() (redis.Conn, error) {
			return redis.DialURL(url)
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}
}

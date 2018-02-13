package pool

import (
	"time"

	"github.com/garyburd/redigo/redis"
	log "github.com/sirupsen/logrus"
)

var (
	pool *redis.Pool
)

func NewRedisPool(server string) {
	if pool == nil {
		pool = newPool(server)
	}
}

func GetRedisPool(server string) *redis.Pool {
	if pool == nil {
		pool = newPool(server)
	}
	return pool
}

func Set_TTL(conn redis.Conn, key, val string, time int) {
	_, err := conn.Do("SET", key, val)
	if err != nil {
		log.WithFields(log.Fields{
			"prefix": "pool.Set_TTL.SET",
			"error":  err.Error(),
		}).Error("error")
	}
	_, err = conn.Do("EXPIRE", key, time)
	if err != nil {
		log.WithFields(log.Fields{
			"prefix": "pool.Set_TTL.EXPIRE",
			"error":  err.Error(),
		}).Error("error")
	}
}

func newPool(server string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     20,
		IdleTimeout: 200 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server,
				redis.DialConnectTimeout(30*time.Second),
				redis.DialReadTimeout(30*time.Second),
				redis.DialWriteTimeout(30*time.Second))
			if err != nil {
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}

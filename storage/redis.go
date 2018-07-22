package storage

import (
	"log"
	"time"

	"github.com/gomodule/redigo/redis"
)

type RedisCache struct {
	pool     *redis.Pool
	cacheTTL int
}

func NewRedisCache(server string, cacheTTL int) *RedisCache {
	rc := RedisCache{
		pool: &redis.Pool{
			MaxActive:   5,
			IdleTimeout: 30 * time.Second,
			Dial: func() (redis.Conn, error) {
				c, err := redis.Dial("tcp", server)
				if err != nil {
					return nil, err
				}
				return c, err
			},
			TestOnBorrow: func(c redis.Conn, t time.Time) error {
				_, err := c.Do("PING")
				return err
			},
		},

		cacheTTL: cacheTTL,
	}

	return &rc
}

func (rc *RedisCache) Set(key string, value []byte) error {
	conn := rc.pool.Get()
	defer conn.Close()

	_, err := conn.Do("SETEX", key, rc.cacheTTL, value)
	if err != nil {
		log.Println("[WARN] error while redis SETEX: ", err.Error())
		return err
	}

	return nil
}

func (rc *RedisCache) Get(key string) ([]byte, error) {
	conn := rc.pool.Get()
	defer conn.Close()

	var data []byte
	data, err := redis.Bytes(conn.Do("GET", key))
	if err != nil {
		log.Println("[WARN] error while redis GET: ", err.Error())
		return nil, err
	}

	return data, err
}

func (rc *RedisCache) Exists(key string) (bool, error) {
	conn := rc.pool.Get()
	defer conn.Close()

	ok, err := redis.Bool(conn.Do("EXISTS", key))
	if err != nil {
		log.Println("[WARN] error while redis EXISTS: ", err.Error())
		return ok, err
	}

	return ok, err
}

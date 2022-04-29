package redis

import (
	"errors"
	"github.com/gomodule/redigo/redis"
	"github.com/miaocansky/go-tool/cache"
	"time"
)

type RedisStore struct {
	pool *redis.Pool
	ttl  time.Duration
	cache.BaseCache
}

func NewRedisCache(network string, address string, db int, password string, defaultExpiration time.Duration) *RedisStore {
	pool := &redis.Pool{
		MaxActive:   512,
		MaxIdle:     10,
		Wait:        false,
		IdleTimeout: 3 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial(
				network,
				address,
				redis.DialDatabase(db),
			)
			if err != nil {
				return nil, err
			}
			if len(password) > 0 { // 有密码的情况
				if _, err := c.Do("AUTH", password); err != nil {
					c.Close()
					return nil, err
				}
			} else { // 没有密码的时候 ping 连接
				if _, err := c.Do("ping"); err != nil {
					c.Close()
					return nil, err
				}
			}
			return c, err
		},
	}
	return &RedisStore{pool: pool, ttl: defaultExpiration}
}

func (c *RedisStore) exec(commandName string, args ...interface{}) (interface{}, error) {
	if len(args) < 1 {
		return nil, errors.New("missing required arguments")
	}
	pool := c.pool.Get()
	defer pool.Close()
	return pool.Do(commandName, args...)
}
func (c *RedisStore) Has(key string) bool {
	i, err := redis.Int(c.exec("Exists", key))
	if err != nil {
		return false
	}
	return i == 1
}

func (c *RedisStore) Del(key string) error {
	_, err := c.exec("Del", key)
	return err
}

func (c *RedisStore) Get(key string) interface{} {
	reply, err := c.exec("Get", key)
	bytes, err := redis.Bytes(reply, err)
	res := c.Unmarshal(bytes, err)
	return res
}
func (c *RedisStore) Set(key string, val interface{}, ttl time.Duration) error {
	bytes, err := c.Marshal(val)
	commandName := "SetEx"
	if err != nil {
		return err
	}
	if ttl == 0 {
		commandName = "Set"
		_, err = c.exec(commandName, key, bytes)
	} else {
		_, err = c.exec(commandName, key, int64(ttl/time.Second), bytes)
	}

	return err

}

func (c *RedisStore) Close() error {
	return c.pool.Close()
}

func (c *RedisStore) GetHash(key, field string) interface{} {
	//reply, err := c.exec("HGet", key, field)
	str, err := redis.String(c.exec("HGet", key, field))
	if err != nil {
		return nil
	}
	return str
}

func (c *RedisStore) SetHash(key string, val interface{}) error {
	flat := redis.Args{}.Add(key).AddFlat(val)
	_, err := c.exec("HSet", flat...)
	return err
}

func (c *RedisStore) SetHNX(key string, val interface{}) error {
	flat := redis.Args{}.Add(key).AddFlat(val)
	_, err := c.exec("Hsetnx", flat...)
	return err
}

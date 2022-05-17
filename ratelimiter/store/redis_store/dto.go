package redis_store

import (
	"github.com/go-redis/redis/v8"
	"sync"
	"time"
)

type BaseLimiter struct {
	sync.Mutex
	scriptSHA1  string
	key         string
	redisClient redis.Cmdable
}

type RedisCounterLimiter struct {
	BaseLimiter
	duration time.Duration
	capacity int
	num      int64
}

type RedisTokenLimiter struct {
	BaseLimiter
	duration time.Duration
	rate     int
	capacity int
	num      int64
}

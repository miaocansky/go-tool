package redis_store

import (
	"context"
	"crypto/sha1"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

func NewRedisTokenLimiter(client redis.Cmdable, key string, duration time.Duration, capacity int, rate int) (*RedisTokenLimiter, error) {

	bgCtx := context.Background()
	_, err := client.Ping(bgCtx).Result()
	if err != nil {
		return nil, err
	}

	script := TokenBucketScript
	scriptSHA1 := fmt.Sprintf("%x", sha1.Sum([]byte(script)))
	limiter := RedisTokenLimiter{
		BaseLimiter: BaseLimiter{
			key:         key,
			scriptSHA1:  scriptSHA1,
			redisClient: client,
		},
		duration: duration,
		rate:     rate,
		capacity: capacity,
	}

	limiter.scriptSHA1 = scriptSHA1

	if !limiter.redisClient.ScriptExists(bgCtx, limiter.scriptSHA1).Val()[0] {
		limiter.redisClient.ScriptLoad(bgCtx, script).Val()
	}

	return &limiter, nil
}
func (r *RedisTokenLimiter) Take() (bool, error) {
	r.Lock()
	defer r.Unlock()
	bgCtx := context.Background()
	// 2. try to get from redis
	x, err := r.redisClient.EvalSha(
		bgCtx,
		r.scriptSHA1,
		[]string{r.key},
		int(r.duration/time.Microsecond),
		r.rate,
		r.capacity).Result()

	if err != nil {
		return false, err
	}
	count := x.(int64)
	if count <= 0 {
		return false, nil
	} else {
		count--
		r.num = count
		return true, nil
	}
}

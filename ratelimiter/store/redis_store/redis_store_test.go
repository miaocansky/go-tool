package redis_store

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/vearne/ratelimit"
	"sync"
	"testing"
	"time"
)

func TestRedisCounterLimiter_Take(t *testing.T) {

	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		//Password: "xxx", // password set
		DB: 11, // use default DB
	})

	limiter, err := NewRedisCounterLimiter(client,
		"key:countliimter",
		1*time.Second/2,
		5,
	)

	if err != nil {
		fmt.Println("error", err)
		return
	}

	//time.Sleep(10 * time.Second)

	var wg sync.WaitGroup
	total := 100
	start := time.Now()
	for i := 0; i < total; i++ {
		go countConsume(limiter, &wg)
	}
	wg.Wait()
	cost := time.Since(start)
	fmt.Println("cost", cost, "rate", float64(total)/cost.Seconds())

}

func countConsume(r ratelimit.Limiter, group *sync.WaitGroup) {
	group.Add(1)
	defer group.Done()

	ok, err := r.Take()
	if err != nil {
		ok = false
		fmt.Println("error", err)
	}
	if ok {
		fmt.Println("成功")
	} else {
		fmt.Println("失败")
		//time.Sleep(time.Duration(rand.Intn(100)+1) * time.Millisecond)
	}
}

func TestRedisTokenLimiter_Take(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		//Password: "xxx", // password set
		DB: 11, // use default DB
	})

	limiter, err := NewRedisTokenLimiter(client,
		"key:tokenliimter",
		10*time.Second,
		3,
		5,
	)

	if err != nil {
		fmt.Println("error", err)
		return
	}
	limiter.Take()
	limiter.Take()
	limiter.Take()
	limiter.Take()
	limiter.Take()
	time.Sleep(5 * time.Second)
	var wg sync.WaitGroup
	total := 100
	start := time.Now()
	for i := 0; i < total; i++ {
		go tokenConsume(limiter, &wg)
	}
	wg.Wait()
	cost := time.Since(start)
	fmt.Println("cost", cost, "rate", float64(total)/cost.Seconds())
}

func tokenConsume(r ratelimit.Limiter, group *sync.WaitGroup) {
	group.Add(1)
	defer group.Done()

	ok, err := r.Take()
	if err != nil {
		ok = false
		fmt.Println("error", err)
	}
	if ok {
		fmt.Println("成功")
	} else {
		fmt.Println("失败")
		//time.Sleep(time.Duration(rand.Intn(100)+1) * time.Millisecond)
	}
}

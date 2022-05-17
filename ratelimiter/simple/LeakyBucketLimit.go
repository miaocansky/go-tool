package simple

import (
	"math"
	"sync"
	"time"
)

type LeakyBucketLimit struct {
	rate       float64 //固定每秒出水速率
	capacity   float64 //桶的容量
	water      float64 //桶中当前水量
	lastLeakMs int64   //桶上次漏水时间戳 ms

	lock sync.Mutex
}

func NewLeakyBucketLimit(rate, capacity float64) *LeakyBucketLimit {
	limit := &LeakyBucketLimit{}
	limit.Set(rate, capacity)
	return limit
}

func (l *LeakyBucketLimit) Take() (bool, error) {
	l.lock.Lock()
	defer l.lock.Unlock()

	now := time.Now().UnixNano() / 1e6
	eclipse := float64((now - l.lastLeakMs)) * l.rate / 1000 //先执行漏水
	l.water = l.water - eclipse                              //计算剩余水量
	l.water = math.Max(0, l.water)                           //桶干了
	l.lastLeakMs = now
	if (l.water + 1) < l.capacity {
		// 尝试加水,并且水还未满
		l.water++
		return true, nil
	} else {
		// 水满，拒绝加水
		return false, nil
	}
}

func (l *LeakyBucketLimit) Set(rate, capacity float64) {
	l.rate = rate
	l.capacity = capacity
	l.water = 0
	l.lastLeakMs = time.Now().UnixNano() / 1e6
}

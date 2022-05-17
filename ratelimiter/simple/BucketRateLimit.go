package simple

import (
	"fmt"
	"sync"
	"time"
)

type BucketLimit struct {
	cycle              time.Duration // 时间间隔
	capacity           int64         //桶的容量
	rate               int64         // 时间间隔放入token数量
	tokens             int64         //  桶token数
	lock               sync.Mutex
	recentPutTokenTime time.Time // 最近放token时间
}

//
//  NewBucketLimit
//  @Description: 创建实体
//  @param capacity
//  @param rate
//  @param cycle
//  @return *BucketLimit
//
func NewBucketLimit(capacity, rate int64, cycle time.Duration) *BucketLimit {

	limit := &BucketLimit{
		capacity:           capacity,
		rate:               rate,
		cycle:              cycle,
		recentPutTokenTime: time.Now(),
	}
	return limit
}

//
//  Take
//  @Description: 获取token
//  @receiver bucketLimit
//  @return bool
//
func (bucketLimit *BucketLimit) Take() (bool, error) {
	bucketLimit.lock.Lock()
	defer bucketLimit.lock.Unlock()
	tokens := bucketLimit.updateToken()
	if tokens <= 0 {
		return false, nil
	}
	return true, nil
}

//
//  updateToken
//  @Description: 更新当前桶的tokens数量
//  @receiver bucketLimit
//
func (bucketLimit *BucketLimit) updateToken() int64 {
	num := float64(time.Now().Sub(bucketLimit.recentPutTokenTime) / bucketLimit.cycle)
	addTokens := int64(num * (float64)(bucketLimit.rate))
	if addTokens > 0 {
		fmt.Println(addTokens)

	}
	if bucketLimit.tokens+addTokens > bucketLimit.capacity {
		bucketLimit.tokens = bucketLimit.capacity
	} else {
		bucketLimit.tokens = bucketLimit.tokens + addTokens
	}
	if addTokens > 0 {
		bucketLimit.recentPutTokenTime = time.Now()
	}
	tokens := bucketLimit.tokens
	if tokens > 0 {
		bucketLimit.tokens--
	}
	return tokens

}

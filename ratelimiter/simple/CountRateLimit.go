package simple

import (
	"sync"
	"time"
)

type CountLimit struct {
	startTime time.Time     // 起始时间
	cycle     time.Duration // 时间间隔
	rate      float64       // 时间间隔内次数
	num       float64       // 请求成功次数
	lock      sync.Mutex
}

//
//  NewCountLimit
//  @Description: 构建实体
//  @param rate
//  @param cycle
//  @return *CountLimit
//
func NewCountLimit(rate float64, cycle time.Duration) *CountLimit {
	startTime := time.Now()
	return &CountLimit{
		startTime: startTime,
		cycle:     cycle,
		rate:      rate,
	}
}

//
//  Take
//  @Description: 获取是否限数
//  @receiver countLimit
//  @return bool
//
func (countLimit *CountLimit) Take() (bool, error) {
	countLimit.lock.Lock()
	defer countLimit.lock.Unlock()
	if countLimit.num == (countLimit.rate - 1) {
		if time.Now().Sub(countLimit.startTime) >= countLimit.cycle {
			countLimit.rest()
			return true, nil

		} else {
			return false, nil
		}

	} else {
		countLimit.num++
		return true, nil
	}
}

//
//  rest
//  @Description: 重置次数和时间
//  @receiver countLimit
//
func (countLimit *CountLimit) rest() {
	countLimit.num = 0
	countLimit.startTime = time.Now()

}

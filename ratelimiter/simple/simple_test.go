package simple

import (
	"log"
	"sync"
	"testing"
	"time"
)

func TestCountLimit(t *testing.T) {
	limit := NewCountLimit(3, time.Second)
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		log.Println("创建请求:", i)
		go func(i int) {
			take, _ := limit.Take()
			if take {
				log.Println("响应请求:", i)
			}
			wg.Done()
		}(i)

		time.Sleep(200 * time.Millisecond)
	}
	wg.Wait()

}

func TestBucketLimit(t *testing.T) {
	//capacity, rate int64, cycle time.Duration
	limit := NewBucketLimit(2, 3, time.Second)
	time.Sleep(1 * time.Second)
	//var wg sync.WaitGroup

	for i := 0; i < 10000; i++ {
		//wg.Add(1)
		//log.Println("创建请求:", i)
		go func(i int) {
			take, _ := limit.Take()
			if take {
				log.Println("响应请求:", i)
			}
			//wg.Done()
		}(i)

		time.Sleep(220 * time.Millisecond)
	}
	//wg.Wait()
	time.Sleep(20 * time.Second)
}

func TestLeakyBucket_Set(t *testing.T) {
	var wg sync.WaitGroup
	limit := NewLeakyBucketLimit(3, 2)
	//time.Sleep(time.Second)
	for i := 0; i < 10; i++ {
		wg.Add(1)
		log.Println("创建请求:", i)
		go func(i int) {
			take, _ := limit.Take()
			if take {
				log.Println("响应请求:", i)
			}
			wg.Done()
		}(i)

		time.Sleep(260 * time.Millisecond)
	}
	wg.Wait()

}

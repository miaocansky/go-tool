package cache

import (
	"fmt"
	"github.com/miaocansky/go-tool/cache/memcached"
	"github.com/miaocansky/go-tool/cache/memory"
	"github.com/miaocansky/go-tool/cache/redis"
	"testing"
	"time"
)

type KeyVal struct {
	Msg  string
	Code int64
}

//
func TestRedisDrive(t *testing.T) {
	//
	k := new(KeyVal)
	k.Msg = "ceshi"
	k.Code = 10

	//m := map[string]string{
	//	"description": "oschina",
	//	"url":         "http://my.oschina.net/myblog",
	//	"author":      "xxbandy",
	//}
	Register("redis", redis.NewRedisCache("tcp", "127.0.0.1:6379", 10, "", time.Second*100))
	SetDefaultName("redis")
	//DefaultDriver().SetHash("testh", m)
	//hash := DefaultDriver().GetHash("testh", "description")
	//s := hash.(string)
	//fmt.Println(s)

	DefaultDriver().Set("aaa", "test", time.Second*100)
	get := DefaultDriver().Get("aaa")
	fmt.Println(get)

}
func TestMemcachedDriver(t *testing.T) {

	k := new(KeyVal)
	k.Msg = "ceshi"
	k.Code = 10
	address := "127.0.0.1:32768"
	Register("memcached", memcached.NewMemcachedStore(address))
	SetDefaultName("memcached")
	DefaultDriver().Set("key", k, time.Second*100)

	get := DefaultDriver().Get("key")

	DefaultDriver().Close()

	fmt.Println(get)
}

func TestGoCacheDriver(t *testing.T) {
	k := new(KeyVal)
	k.Msg = "ceshi"
	k.Code = 10
	Register("goCache", memory.NewMemoryStore())
	SetDefaultName("goCache")
	//DefaultDriver().Set("key", k, time.Minute*30)
	get := DefaultDriver().Get("key")
	fmt.Println(get)
	DefaultDriver().Close()

}

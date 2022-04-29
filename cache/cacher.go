package cache

import (
	"io"
	"time"
)

type SimpleCacher interface {
	io.Closer
	Has(key string) bool
	Del(key string) error
	Get(key string) interface{}
	Set(key string, val interface{}, ttl time.Duration) error
	//SetHash(key string, val interface{}) error
	//GetHash(key, field string) interface{}
	//
	//GetMulti(keys []string) map[string]interface{}
	//SetMulti(values map[string]interface{}, ttl time.Duration) error
	//DelMulti(keys []string) error
}

type RedisCacher interface {
	SimpleCacher
	SetHash(key string, val interface{}) error
	GetHash(key, field string) interface{}
}

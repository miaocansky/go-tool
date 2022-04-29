package memory

import (
	"github.com/miaocansky/go-tool/cache"
	GoCache "github.com/patrickmn/go-cache"
	"os"
	"time"
)

type MemoryStore struct {
	cache.BaseCache
	goC *GoCache.Cache
}

func NewMemoryStore() *MemoryStore {
	m := &MemoryStore{
		goC: GoCache.New(3*time.Minute, time.Minute*10),
	}
	//path := GetDataPath("go-cache-cache.dat")
	//exists, _ := PathExists(path)
	//if exists {
	//	m.goC.LoadFile(path)
	//}
	return m
}

func (m *MemoryStore) Close() error {
	//path := GetDataPath("go-cache-cache.dat")
	//err := m.goC.SaveFile(path)
	//fmt.Println(err)
	m.goC = nil
	return nil
}

func (m *MemoryStore) Has(key string) (b bool) {
	_, b = m.goC.Get(key)
	return
}

func (m *MemoryStore) Del(key string) error {
	m.goC.Delete(key)
	return nil
}

func (m *MemoryStore) Get(key string) interface{} {
	val, _ := m.goC.Get(key)
	return val
}

func (m *MemoryStore) Set(key string, val interface{}, ttl time.Duration) error {
	m.goC.Set(key, val, ttl)
	return nil
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	//当为空文件或文件夹存在
	if err == nil {
		return true, nil
	}
	//os.IsNotExist(err)为true，文件或文件夹不存在
	if os.IsNotExist(err) {
		return false, nil
	}
	//其它类型，不确定是否存在
	return false, err
}
func GetDataPath(name string) string {
	path, _ := os.Getwd()
	path = path + "/" + name
	return path

}

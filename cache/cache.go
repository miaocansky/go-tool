package cache

import (
	"encoding/json"
	"fmt"
)

var std = NewManager()

func Register(driverName string, driver SimpleCacher) *Manager {
	std.Register(driverName, driver)
	return std
}
func UnRegister(driverName string) bool {
	ok := std.UnRegister(driverName)
	return ok
}

func SetDefaultName(driverName string) {
	std.SetDefaultName(driverName)
}
func DefaultDriver() SimpleCacher {
	return std.DefaultDriver()
}
func Use(driverName string) SimpleCacher {
	return std.Use(driverName)
}

func GetDriver(driverName string) SimpleCacher {
	return std.GetDriver(driverName)
}
func CloseAll() error {
	return std.CloseAll()
}

type BaseCache struct {
} //
//  Unmarshal
//  @Description: 解析数据 json
//  @receiver c
//  @param bytes
//  @param err
//  @return interface{}
//
func (c *BaseCache) Unmarshal(bytes []byte, err error) interface{} {
	if err != nil {
		// log err
		fmt.Errorf("error: %T", err)
		return nil
	}
	var res interface{}

	err = json.Unmarshal(bytes, &res)
	if err != nil {
		// log err
		fmt.Errorf("error: %T", err)
	}
	return res
}

//
//  Marshal
//  @Description: 生成 json
//  @receiver c
//  @param val
//  @return []byte
//  @return error
//
func (c *BaseCache) Marshal(val interface{}) ([]byte, error) {
	return json.Marshal(val)
}

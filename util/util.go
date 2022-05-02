package util

import (
	"github.com/fatih/structs"
	"github.com/goinggo/mapstructure"
	"reflect"
)

//
//  StructAssign
//  @Description: 结构体绑定
//  @param binding 空结构体
//  @param value 有数据的结构体
//
func StructAssign(binding interface{}, value interface{}) {
	bVal := reflect.ValueOf(binding).Elem() //获取reflect.Type类型
	vVal := reflect.ValueOf(value).Elem()   //获取reflect.Type类型
	vTypeOfT := vVal.Type()
	for i := 0; i < vVal.NumField(); i++ {
		// 在要修改的结构体中查询有数据结构体中相同属性的字段，有则修改其值
		name := vTypeOfT.Field(i).Name
		if ok := bVal.FieldByName(name).IsValid(); ok {
			bVal.FieldByName(name).Set(reflect.ValueOf(vVal.Field(i).Interface()))
		}
	}
}

//
//  StructToMap
//  @Description: struct转map
//  @param s
//  @return map[string]interface{}
//
func StructToMap(s interface{}) map[string]interface{} {
	out := make(map[string]interface{})
	m := structs.Map(s)
	for key, val := range m {
		switch v := val.(type) {
		case map[string]interface{}:
			out = MergeMap(out, v)
		default:
			out[key] = v
		}
	}
	return out
}

//
//  MergeMap
//  @Description: 合并map
//  @param mObj
//  @return map[string]interface{}
//
func MergeMap(mObj ...map[string]interface{}) map[string]interface{} {
	newObj := make(map[string]interface{})
	for _, m := range mObj {
		for k, v := range m {
			newObj[k] = v
		}
	}
	return newObj
}

//
//  MapToStruct
//  @Description:  map 转 stuct
//  @param m  map
//  @param st  struct
//  @return error
//
func MapToStruct(m interface{}, st interface{}) error {
	return mapstructure.Decode(m, st)
}

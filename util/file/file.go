package file

import "os"

//
//  GetPath
//  @Description: 获取当前地址
//  @return string
//  @return error
//
func GetPath() (string, error) {
	return os.Getwd()
}

//
//  PathExists
//  @Description: 文件地址是否存在
//  @param path
//  @return bool
//  @return error
//
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

//
//  MakeDir
//  @Description: 创建文件目录
//  @param src
//  @return error
//
func MakeDir(src string) error {
	err := os.MkdirAll(src, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

package upload

import (
	"github.com/miaocansky/go-tool/util/file"
	"io/ioutil"
	"mime/multipart"
)

type UploadUtil struct {
	Path string //项目当前相对地址
}

func NewUploadUtil() *UploadUtil {
	return &UploadUtil{}
}
func NewPathUploadUtil(path string) *UploadUtil {
	return &UploadUtil{
		Path: path,
	}

}

//
//  Upload
//  @Description:文件上传
//  @param file 文件流
//  @param name 文件名称
//  @return string  文件地址
//  @return string url地址
//  @return error  异常
//
func (uploadUtil *UploadUtil) Upload(file multipart.File, name string) (string, string, error) {
	savePath, path, err := uploadUtil.GetPath()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return savePath, "", err
	}
	err = ioutil.WriteFile(savePath+name, data, 0666)
	if err != nil {
		return savePath, "", err
	}
	return savePath + name, path + name, nil
}

func (uploadUtil *UploadUtil) GetPath() (string, string, error) {

	var savePath string
	var path string
	DirPath, err := file.GetPath()
	if err != nil {
		return "", "", err
	}
	if uploadUtil.Path != "" && len(uploadUtil.Path) > 0 {
		//savePath = uploadUtil.SavePath
		path = uploadUtil.Path
		savePath = DirPath + "/" + path
	} else {

		path = "upload/"
		savePath = DirPath + "/" + path
	}

	exists, _ := file.PathExists(savePath)
	if !exists {
		//  创建目录
		err := file.MakeDir(savePath)
		if err != nil {
			return "", "", err
		}
	}

	return savePath, path, nil
}

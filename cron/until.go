package cron

import (
	"errors"
	"io/ioutil"
	"net/http"
)

func HttpRequest(method string, url string) (string, error) {
	//声明client 参数为默认
	client := &http.Client{}
	reqest, err := http.NewRequest(method, url, nil)
	if err != nil {
		return "", err
	}
	//处理返回结果
	response, err := client.Do(reqest)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	if response.StatusCode == 200 {
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return "", err
		} else {
			return string(body), nil
		}
	} else {
		return "", errors.New("http请求状态异常")
	}
}

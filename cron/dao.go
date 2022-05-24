package cron

type CronTask struct {
	Id         int64  // id
	Name       string // 名称
	Spec       string // crontab 表达式
	Command    string // 命令
	Type       int64  // 运行类型1:shell 2:http
	HttpMethod int64  // http请求方式1:get 2:post
}

package cron

type CronTask struct {
	Id         int64       // id
	Name       string      // 名称
	Spec       string      // crontab 表达式
	Command    string      // 命令
	Type       int64       // 运行类型1:shell 2:http
	HttpMethod int64       // http请求方式1:get 2:post
	ExecNum    int64       // 当前执行的任务个数
	IsSingle   bool        // 是否只允许单个执行
	DelayClose bool        // 是否延迟关闭
	CallBack   CallBackFuc //回调处理
}

type ResultExecData struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
}

type CallBackFuc func(data ResultExecData, task CronTask)

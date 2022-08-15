package cron

import (
	"github.com/jakecoffman/cron"
	"github.com/miaocansky/go-tool/log"
	"github.com/miaocansky/go-tool/log/zap"
	"strconv"
	"sync"
)

type Croner interface {
	// 任务id（key）
	getJobId(id int64) string
	// 开始
	Start()
	// 停止
	Stop()
	// 添加任务
	AddTask(task *CronTask)
	// RemoveTask 删除任务如果该任务没有在执行中 则返回true并且删除任务  如果执行中则返回false 进行延时删除任务（等该任务都已经执行结束以后 删除改任务 并且该任务不会再添加执行 ）
	RemoveTask(taskId int64) bool
	// 添加工作
	AddJob(task *CronTask) cron.FuncJob
	// 外放执行脚本
	ExecuteJob(task *CronTask)
	// 获取所有的任务Ids
	GetAllTaskIds() []int64
	// 获取所有的任务数据包含每个执行状态
	GetAllTasksDetailList() map[string]CronTask
	// 获取cron
	GetCron() *cron.Cron
}

type CronServer struct {
	taskMaps  sync.Map
	cronLog   log.Logger
	debug     bool
	pfx       string
	cron      *cron.Cron
	taskCount *taskCount
}

func NewCronComServer() *CronServer {
	return NewCronServer(false)
}

func NewCronDebugServer() *CronServer {
	return NewCronServer(true)
}

func NewCronServer(debug bool) *CronServer {

	cronServer := &CronServer{
		debug: debug,
		cron:  cron.New(),
		taskCount: &taskCount{
			wg:   sync.WaitGroup{},
			exit: make(chan struct{}),
		},
	}
	if debug {
		config := zap.ZapConfig{
			Director:      "log/cron",
			Level:         "info",
			ShowLine:      false,
			StacktraceKey: "stacktrace",
			EncodeLevel:   "LowercaseColorLevelEncoder",
			Format:        "console",
			Prefix:        "[cron ]",
			LinkName:      "latest_log",
			LogInConsole:  true,
		}
		log.Register("zap", zap.NewZapUtil(config))
		cronLog := log.DefaultDriver()
		cronServer.cronLog = cronLog
	}

	return cronServer

}

func (cS *CronServer) getJobId(id int64) string {
	idStr := strconv.FormatInt(id, 10)
	return cS.pfx + idStr

}
func (cS *CronServer) GetCron() *cron.Cron {
	return cS.cron
}

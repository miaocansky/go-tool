package cron

import (
	"github.com/jakecoffman/cron"
	"github.com/miaocansky/go-tool/log"
	"github.com/miaocansky/go-tool/log/zap"
	"strconv"
	"sync"
)

type Croner interface {
	getJobId(id int64) string

	// 开始
	Start()
	// 停止
	Stop()

	// 添加任务
	AddTask(task *CronTask)

	// RemoveTask 删除定时任务
	RemoveTask(taskId int64)

	// 添加工作
	AddJob(task *CronTask) cron.FuncJob
	// 外放执行脚本
	ExecuteJob(task *CronTask)

	GetAllTasks()
}

type CronServer struct {
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

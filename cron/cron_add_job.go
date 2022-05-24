package cron

import (
	"fmt"
	"github.com/jakecoffman/cron"
)

func (cS *CronServer) AddJob(task *CronTask) cron.FuncJob {
	return func() {
		cS.taskCount.Add()
		defer cS.taskCount.Done()
		if cS.debug && cS.cronLog != nil {
			msg := fmt.Sprintf("执行任务：(%d)%s [%s]", task.Id, task.Name, task.Spec)
			cS.cronLog.Info(msg)
		}

		//fmt.Println(msg)
	}

}

package cron

import (
	"strings"
)

func (cS *CronServer) AddTask(task *CronTask) {
	spec := strings.TrimSpace(task.Spec)
	key := cS.getJobId(task.Id)
	cS.cron.AddFunc(spec, cS.AddJob(task), key)

	_, ok := cS.taskMaps.Load(key)
	if !ok {
		// 不存在
		cS.taskMaps.Store(key, task)
	}
}

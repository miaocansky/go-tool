package cron

import "strings"

func (cS *CronServer) AddTask(task *CronTask) {
	spec := strings.TrimSpace(task.Spec)
	name := cS.getJobId(task.Id)
	cS.cron.AddFunc(spec, cS.AddJob(task), name)

}

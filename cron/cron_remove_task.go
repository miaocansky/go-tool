package cron

func (cS *CronServer) RemoveTask(taskId int64) {
	name := cS.getJobId(taskId)
	cS.cron.RemoveJob(name)
}

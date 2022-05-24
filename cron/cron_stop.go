package cron

func (cS *CronServer) Stop() {
	cS.cron.Stop()
	cS.taskCount.Exit()
}

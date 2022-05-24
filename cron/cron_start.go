package cron

func (cS *CronServer) Start() {
	cS.cron.Start()
	cS.taskCount.Wait()

}

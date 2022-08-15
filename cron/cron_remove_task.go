package cron

/**
 * 删除任务如果该任务没有在执行中 则返回true 如果执行中则返回false 进行延时删除任务（等该任务都已经执行结束以后 删除改任务 并且该任务不会再添加执行 ）
 */
func (cS *CronServer) RemoveTask(taskId int64) bool {
	key := cS.getJobId(taskId)
	task, ok := cS.taskMaps.Load(key)
	var isDelete bool = false
	if ok {
		cronTask := task.(*CronTask)
		if cronTask.ExecNum == 0 {
			// 如果未在执行可以删除
			isDelete = true
		} else {
			//如果执行中则设置延时关闭
			cS.SetDelayClose(cronTask)
		}
	} else {
		isDelete = true
	}
	if isDelete {
		cS.cron.RemoveJob(key)
		cS.taskMaps.Delete(key)

	}
	return isDelete

}

/*
 * 设置任务延迟删除任务
 */
func (cS *CronServer) SetDelayClose(task *CronTask) {
	key := cS.getJobId(task.Id)
	orgTask, ok := cS.taskMaps.Load(key)
	if ok {
		orgCronTask := orgTask.(*CronTask)
		orgCronTask.DelayClose = true
		cS.taskMaps.Store(key, orgCronTask)
	}
}

/**
 * 验证任务延迟删除任务
 */
func (cS *CronServer) checkedDelayClose(task *CronTask) {
	key := cS.getJobId(task.Id)
	orgTask, ok := cS.taskMaps.Load(key)
	if ok {
		orgCronTask := orgTask.(*CronTask)
		if orgCronTask.DelayClose && orgCronTask.ExecNum == 0 {
			cS.RemoveTask(task.Id)
		}
	}

}

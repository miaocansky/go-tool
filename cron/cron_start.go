package cron

import "strconv"

func (cS *CronServer) Start() {
	cS.cron.Start()
	cS.taskCount.Wait()
}

func (cS *CronServer) GetAllTaskIds() []int64 {
	entries := cS.cron.Entries()
	tasks := make([]int64, 0)
	for _, entry := range entries {
		idStr := entry.Name
		id, _ := strconv.ParseInt(idStr, 10, 64)
		tasks = append(tasks, id)
	}
	return tasks

}
func (cS *CronServer) GetAllTasksDetailList() map[string]CronTask {
	tasks := make(map[string]CronTask)
	cS.taskMaps.Range(func(k, v interface{}) bool {
		key := k.(string)
		task := v.(*CronTask)
		tasks[key] = *task
		return true
	})
	return tasks
}
func (cS *CronServer) GetTaskDetail(id int64) *CronTask {
	key := cS.getJobId(id)
	val, ok := cS.taskMaps.Load(key)
	if !ok {
		return nil
	}
	task := val.(*CronTask)
	return task
}

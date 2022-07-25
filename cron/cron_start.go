package cron

import "strconv"

func (cS *CronServer) Start() {
	cS.cron.Start()
	cS.taskCount.Wait()
}

func (cS *CronServer) GetAllTasks() []int64 {
	entries := cS.cron.Entries()
	tasks := make([]int64, 0)
	for _, entry := range entries {
		idStr := entry.Name
		id, _ := strconv.ParseInt(idStr, 10, 64)
		tasks = append(tasks, id)
	}
	return tasks

}

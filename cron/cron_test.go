package cron

import (
	"testing"
)

func TestNewCronServer(t *testing.T) {

	cs := NewCronDebugServer()
	spec1 := "*/3 * * * * ?"
	task := &CronTask{
		Id:         1,
		Name:       "测试任务",
		Spec:       spec1,
		Command:    "command",
		Type:       0,
		HttpMethod: 0,
	}

	spec2 := "*/6 * * * * ?"
	task2 := &CronTask{
		Id:         2,
		Name:       "测试任务2",
		Spec:       spec2,
		Command:    "command",
		Type:       0,
		HttpMethod: 0,
	}

	cs.AddTask(task)
	cs.AddTask(task2)
	cs.Start()

}

package cron

import (
	"fmt"
	"testing"
	"time"
)

func TestNewCronServer(t *testing.T) {

	cs := NewCronDebugServer()
	spec1 := "*/3 * * * * ?"
	task := &CronTask{
		Id:         1,
		Name:       "测试任务",
		Spec:       spec1,
		Command:    "whoami",
		Type:       1,
		HttpMethod: 0,
		CallBack:   nil,
	}

	spec2 := "*/6 * * * * ?"
	task2 := &CronTask{
		Id:         2,
		Name:       "测试任务2",
		Spec:       spec2,
		Command:    "http://127.0.0.1:8081/public/test/test",
		Type:       2,
		HttpMethod: 1,
		CallBack:   DoSoothing(),
	}

	cs.AddTask(task)
	cs.AddTask(task2)
	cs.Start()

}

func TestCronServer_ExecuteJob(t *testing.T) {
	cs := NewCronDebugServer()
	spec2 := "*/6 * * * * ?"
	task2 := &CronTask{
		Id:         2,
		Name:       "测试任务2",
		Spec:       spec2,
		Command:    "http://127.0.0.1:8081/public/test/test",
		Type:       2,
		HttpMethod: 1,
		CallBack:   DoSoothing(),
	}

	cs.ExecuteJob(task2)

	time.Sleep(5 * time.Second)

}

func DoSoothing() CallBackFuc {

	return func(data ResultExecData, task CronTask) {
		fmt.Println(task.Name)
		fmt.Println(data.Msg)
	}
}

package cron

import "sync"

type taskCount struct {
	wg   sync.WaitGroup
	exit chan struct{}
}

func (taskCount *taskCount) Add() {
	taskCount.wg.Add(1)

}
func (taskCount *taskCount) Done() {
	taskCount.wg.Done()
}

func (taskCount *taskCount) Exit() {
	taskCount.wg.Done()
	<-taskCount.exit
}
func (taskCount *taskCount) Wait() {
	taskCount.Add()
	taskCount.wg.Wait()
	close(taskCount.exit)
}

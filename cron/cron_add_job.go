package cron

import (
	"encoding/json"
	"fmt"
	"github.com/jakecoffman/cron"
	"os/exec"
)

/*
 * 设置任务开始执行
 */
func (cS *CronServer) SetTaskExecStart(task *CronTask) {
	key := cS.getJobId(task.Id)
	orgTask, ok := cS.taskMaps.Load(key)
	if ok {
		orgCronTask := orgTask.(*CronTask)
		//orgCronTask.IsExec = true
		orgCronTask.ExecNum = orgCronTask.ExecNum + 1
		cS.taskMaps.Store(key, orgCronTask)
	}
}

/*
 * 设置任务开始完成
 */
func (cS *CronServer) SetTaskExecFinish(task *CronTask) {
	key := cS.getJobId(task.Id)
	orgTask, ok := cS.taskMaps.Load(key)
	if ok {
		orgCronTask := orgTask.(*CronTask)
		if orgCronTask.ExecNum > 0 {
			orgCronTask.ExecNum = orgCronTask.ExecNum - 1
		}

		cS.taskMaps.Store(key, orgCronTask)
	}
}

/*
 * 任务读取
 */
func (cS *CronServer) GetNowTask(task *CronTask) *CronTask {
	key := cS.getJobId(task.Id)
	orgTask, ok := cS.taskMaps.Load(key)
	if ok {
		orgCronTask := orgTask.(*CronTask)
		return orgCronTask
	}
	return nil
}

func (cS *CronServer) AddJob(oldTask *CronTask) cron.FuncJob {
	return func() {

		task := cS.GetNowTask(oldTask)
		//  如果已经在关闭中的任务 无需在新增任务
		if task.DelayClose {
			return
		}
		// 单个任务并且已经在执行了 那不允许启动第二个任务
		if task.ExecNum > 0 && task.IsSingle {
			return
		}
		cS.SetTaskExecStart(task)
		cS.taskCount.Add()

		defer func() {
			cS.SetTaskExecFinish(task)
			cS.taskCount.Done()
			cS.checkedDelayClose(task)

		}()
		msg := fmt.Sprintf("执行任务：(%d)%s [%s]", task.Id, task.Name, task.Spec)
		result := execute(task)
		if cS.debug && cS.cronLog != nil {
			cS.cronLog.Info(msg)
			resultMsg := fmt.Sprintf("%s [返回状态: %s 放回说明: %s]", msg, result.Code, result.Msg)
			cS.cronLog.Info(resultMsg)
		}
		if task.CallBack != nil {
			go task.CallBack(result, *task)
		}
	}

}
func (cS *CronServer) ExecuteJob(task *CronTask) {

	msg := fmt.Sprintf("执行任务：(%d)%s [%s]", task.Id, task.Name, task.Spec)
	result := execute(task)
	if cS.debug && cS.cronLog != nil {
		//cS.cronLog.Info(msg)
		resultMsg := fmt.Sprintf("%s [返回状态: %s 放回说明: %s]", msg, result.Code, result.Msg)
		cS.cronLog.Info(resultMsg)
	}
	if task.CallBack != nil {
		go task.CallBack(result, *task)
	}

}

func execute(task *CronTask) (result ResultExecData) {
	if task.Type == 1 {
		result = executeShell(task)
	} else if task.Type == 2 {
		result = executeUrl(task)
	} else {
		result.Code = "err"
		result.Msg = "执行类型不存在"
	}
	return result
}

func executeShell(task *CronTask) ResultExecData {
	var result ResultExecData
	cmd := exec.Command(task.Command)
	byteArr, err := cmd.Output()
	if err != nil {
		result.Code = "err"
		result.Msg = "执行异常:" + err.Error()
		return result

	}
	jsonStr := string(byteArr)
	return analysisResData(jsonStr)
}

func executeUrl(task *CronTask) ResultExecData {
	var result ResultExecData
	method := "GET"
	if task.HttpMethod == 2 {
		method = "POST"
	}
	body, err := HttpRequest(method, task.Command)
	if err != nil {
		result.Code = "err"
		result.Msg = err.Error()
		return result
	}
	return analysisResData(body)

}

func analysisResData(jsonStr string) ResultExecData {
	var result ResultExecData
	err := json.Unmarshal([]byte(jsonStr), &result)
	if err != nil {
		result.Code = "err"
		result.Msg = "解析返回值异常:" + err.Error()
		return result
	}
	return result
}

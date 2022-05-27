package cron

import (
	"encoding/json"
	"fmt"
	"github.com/jakecoffman/cron"
	"os/exec"
)

func (cS *CronServer) AddJob(task *CronTask) cron.FuncJob {
	return func() {
		cS.taskCount.Add()
		defer cS.taskCount.Done()
		msg := fmt.Sprintf("执行任务：(%d)%s [%s]", task.Id, task.Name, task.Spec)
		if cS.debug && cS.cronLog != nil {
			//cS.cronLog.Info(msg)
		}
		result := execute(task)
		resultMsg := fmt.Sprintf("%s [返回状态: %s 放回说明: %s]", msg, result.Code, result.Msg)
		cS.cronLog.Info(resultMsg)
		if task.CallBack != nil {
			task.CallBack(result)
		}

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

package main

import "os"

func main() {
	err := os.Mkdir(AppstatusPath, 0777)
	if !os.IsExist(err) {
		os.Exit(10)
	}
	err = os.Mkdir(LogPath, 0777)
	if !os.IsExist(err) {
		WriteError2Appstatus(err.Error(), 9)
	}
	//网络检查
	netCheckFlag, err := NetCheck()
	if err != nil || !netCheckFlag{
		ConnectFail()
		WriteError2Appstatus(err.Error(), 2)
	}
	//任务开始
	TaskStart()
	//读取配置
	GetConf()
	tasks, err := GetTaskConfig()
	if err != nil {
		GetConfFail()
		WriteError2Appstatus(err.Error(), 3)
	}
	GetConfSuccess()
	//任务执行
	TaskRun()
	err = ControlDNSQueryRoutine(tasks)
	if err != nil {
		TaskRunFail()
		WriteError2Appstatus(err.Error(), 1)
	}
	ControlCompareRoutine(tasks)
	TaskRunSuccess()
	//写结果
	WriteResult()
	err = ControlWriteResultRoutine(tasks)
	if err != nil {
		WriteResultFail()
		WriteError2Appstatus(err.Error(), 1)
	}
	WriteResultSuccess()
	//写状态文件
	WriteSuccess2Appstatus()
}

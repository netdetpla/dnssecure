package main

import (
	"errors"
	"fmt"
	"github.com/op/go-logging"
	"os"
	"runtime/debug"
	"time"
)

// 日志
var log = logging.MustGetLogger("example")
var format = logging.MustStringFormatter(
	`%{color}%{time:15:04:05} %{shortfile} %{shortfunc} ▶ %{level:.4s} %{color:reset}  %{message}`,
)

func init() {
	// 日志初始化配置
	backend1 := logging.NewLogBackend(os.Stderr, "", 0)
	backend2 := logging.NewLogBackend(os.Stderr, "", 0)
	backend2Formatter := logging.NewBackendFormatter(backend2, format)
	backend1Leveled := logging.AddModuleLevel(backend1)
	backend1Leveled.SetLevel(logging.ERROR, "")
	logging.SetBackend(backend1Leveled, backend2Formatter)
}

func main() {
	log.Info("start: "+ time.Now().Format("2006/01/02 15:04:05"))
    defer func() {
        if err := recover(); err != nil {
			fmt.Println(string(debug.Stack())[:])
			var r error
			switch x := err.(type) {
			case string: 
				r = errors.New(x)
			case error:
				r = x
			default:
				r = errors.New("")
			}
			log.Error(r.Error())
			log.Error(string(debug.Stack())[:])
        }
    }()
	err := os.Mkdir(AppstatusPath, 0777)
	if err != nil && !os.IsExist(err) {
		fmt.Println(err.Error())
		os.Exit(10)
	}
	defer func() {
		if err := recover(); err != nil {
			WriteError2Appstatus(string(debug.Stack())[:], 16)
		}
	}()
	err = os.Mkdir(LogPath, 0777)
	if err != nil && !os.IsExist(err) {
		WriteError2Appstatus(err.Error(), 9)
	}
	//网络检查
	netCheckFlag := NetCheck()
	if netCheckFlag == 0 {
		ConnectFail()
		WriteError2Appstatus("Can not connect to the Internet.", 22)
	} else if netCheckFlag == -1 {
		ConnectFail()
       WriteError2Appstatus("Ping check timeout.", 21)
	}
	//任务开始
	TaskStart()
	//读取配置
	GetConf()
	tasks, err := GetTaskConfig()
	if err != nil {
		GetConfFail()
		WriteError2Appstatus(err.Error(), 13)
	}
	GetConfSuccess()
	//任务执行
	TaskRun()
	err = ControlDNSQueryRoutine(tasks)
	if err != nil {
		TaskRunFail()
		WriteError2Appstatus(err.Error(), 11)
	}
	TaskRunSuccess()
	//进度
	err = SendProcess(tasks.taskID, tasks.uuid, "DomainInfo", len(tasks.records), true)
	if err != nil {
		WriteResultFail()
		WriteError2Appstatus(err.Error(), 14)
	}
	//写结果
	WriteResult()
	err = ControlWriteResultRoutine(tasks)
	if err != nil {
		WriteResultFail()
		WriteError2Appstatus(err.Error(), 15)
	}
	WriteResultSuccess()
	WriteSuccess2Appstatus()
}

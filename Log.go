package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

const (
	LogPath       = "/tmp/log/"
	AppstatusPath = "/tmp/appstatus/"
	Info          = "[INFO][%s]%s"
	Error         = "[ERROR][%s]%s"
	Warning       = "[WARNING][%s]%s"
)
func GetTime() string{
	return time.Now().Format("2006/01/02 15:04:05")
}

func InfoLog(log string) {
	fmt.Println(Info, GetTime(), log)
}

func ErrorLog(log string) {
	fmt.Println(Error, GetTime(), log)
}

func WarningLog(log string)  {
	fmt.Println(Warning, GetTime(), log)
}

func CreateLogFile(logName string) {
	now := string(time.Now().Unix())
	err := ioutil.WriteFile(LogPath+now+logName, []byte(""), 0644)
	if err != nil {
		os.Exit(10)
	}
}

func TaskStart() {
	CreateLogFile("-1100.log")
}

func GetConf() {
	CreateLogFile("-1200.log")
}

func GetConfSuccess() {
	CreateLogFile("-1202.log")
}

func GetConfFail() {
	CreateLogFile("-1201.log")
}

func TaskRun() {
	CreateLogFile("-1300.log")
}

func TaskRunSuccess() {
	CreateLogFile("-1301.log")
}

func TaskRunFail() {
	CreateLogFile("-1302.log")
}

func WriteResult() {
	CreateLogFile("-1400.log")
}

func WriteResultSuccess() {
	CreateLogFile("-1401.log")
}

func WriteResultFail() {
	CreateLogFile("-1402.log")
}

func TaskSuccess() {
	CreateLogFile("-1102.log")
}

func TaskFail() {
	CreateLogFile("-1101.log")
}

func ConnectFail() {
	CreateLogFile("-1111.log")
}

func WriteSuccess2Appstatus() {
	TaskSuccess()
	err := ioutil.WriteFile(AppstatusPath+"0", []byte(""), 0644)
	if err != nil {
		os.Exit(9)
	}
}

func WriteError2Appstatus(errorInfo string, errorCode int) {
	fmt.Println(errorInfo)
	TaskFail()
	err := ioutil.WriteFile(AppstatusPath+"1", []byte(errorInfo), 0644)
	if err != nil {
		os.Exit(9)
	}
	os.Exit(errorCode)
}

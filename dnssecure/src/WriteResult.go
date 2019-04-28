package main

import (
	"io/ioutil"
	"os"
	"runtime"
	"strings"
	"time"
	"strconv"
)

const ResultPath  = "/tmp/result/"
var resultLine = make(chan string)

func GenerateResultLine(record *Record, taskID string, taskName string) {
	var resultList []string
	detectAsStr := strings.Join(record.detectAs, "+")
	detectCNamesStr := strings.Join(record.detectCNames, "+")
	now := strconv.FormatInt(time.Now().Unix(), 10)
	resultList = append(resultList,	taskID, taskName, record.domain, record.reServer,
		detectAsStr + "/" + detectCNamesStr, now + "\n")
	resultStr := strings.Join(resultList, ";")
	resultLine <- resultStr
	return
}

func ControlWriteResultRoutine(tasks *Task) (err error){
	runtime.GOMAXPROCS(runtime.NumCPU())
	for _, record := range tasks.records {
		go GenerateResultLine(record, tasks.taskID, tasks.taskName)
	}
	var resultContent string
	for i := 0; i < len(tasks.records); i++ {
		resultContent += <- resultLine
	}
	close(resultLine)
	err = os.Mkdir(ResultPath, 0777)
	if err != nil && !os.IsExist(err) {
		return
	}
    totalNum := len(tasks.records)
	err = ioutil.WriteFile(ResultPath + tasks.taskID + ".result",
        []byte(tasks.taskID + "|" + strconv.Itoa(totalNum) + "|" + tasks.subID + "\n" + resultContent), 0644)
	return
}

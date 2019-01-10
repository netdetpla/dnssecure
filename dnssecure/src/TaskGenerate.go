package main

import (
	"encoding/base64"
	"io/ioutil"
	"strings"
	"fmt"
)

const ConfPath = "/tmp/conf/busi.conf"

type Record struct {
	rightRecord *RightRecord
	reServer string
	detectAs []string
	detectCNames []string
	timeoutFlag bool
	result string
	compareType string
}

type Task struct {
	taskID string
	taskName string
	uuid string
	subID string
	records []*Record
}

func GetTaskConfig() (task *Task, err error) {
	task = new(Task)
	taskConfigBase64, err := ioutil.ReadFile(ConfPath)
	if err != nil {
		fmt.Println(err.Error())
        return nil, err
    }
	taskConfigB, err := base64.StdEncoding.DecodeString(string(taskConfigBase64))
	if err != nil {
		fmt.Println(err.Error())
        return nil, err
    }
	taskConfig := strings.Split(string(taskConfigB), ",")
	fmt.Println(taskConfig)

	task.taskID = taskConfig[0]

	//组合域名、递归服务器、正确值
	domains := strings.Split(taskConfig[1], "+")
	rightRecords, err := getRightValue(domains)
	if err != nil {
		fmt.Println(err.Error())
        return nil, err
    }
	reServers := strings.Split(taskConfig[2], "+")
	for _, reServer := range reServers {
		if len(reServer) == 0 {
			continue
		}
		for _, rightRecord := range rightRecords {
			record := new(Record)
			record.rightRecord = rightRecord
			record.reServer = reServer
			task.records = append(task.records, record)
		}
	}

	task.taskName = taskConfig[4]
	task.uuid = taskConfig[5]
	task.subID = taskConfig[6]
	return
}

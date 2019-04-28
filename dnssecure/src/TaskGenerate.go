package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"strings"
)

const ConfPath = "/tmp/conf/busi.conf"

type Record struct {
	domain       string
	reServer     string
	detectAs     []string
	detectCNames []string
	timeoutFlag  bool
}

type Task struct {
	taskID   string
	taskName string
	uuid     string
	subID    string
	records  []*Record
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
	reServers := strings.Split(taskConfig[2], "+")
	for _, reServer := range reServers {
		if len(reServer) == 0 {
			continue
		}
		for _, domain := range domains {
			record := new(Record)
			record.domain = domain
			record.reServer = reServer
			task.records = append(task.records, record)
		}
	}
	task.taskName = taskConfig[4]
	task.uuid = taskConfig[5]
	task.subID = taskConfig[6]
	return
}

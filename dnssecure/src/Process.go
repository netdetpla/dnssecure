package main

import (
	"net"
	"strconv"
)

func SendProcess(taskID string, uuid string, imageType string, count int, finalFlag bool) (err error){
	var processJson = make(map[string]interface{})
	processJson["taskid"] = taskID
	processJson["addnum"] = count
	processJson["id_uuid"] = uuid
	processJson["image_type"] = imageType
	processJson["final"] = finalFlag
	var flag = "true"
	if finalFlag == true {
		flag = "true"
	}else{
		flag = "false"
	}
	var json = "{"+"\"taskid\":\""+taskID+"\",\"addnum\":"+strconv.Itoa(count)+",\"id_uuid\":\""+uuid+"\",\"image_type\":\""+imageType+"\",\"final\":"+flag+"}"
	conn, err := net.Dial("udp", "10.96.129.4:7011")
    if err != nil {
        return
    }
    _, err = conn.Write([]byte(json))
    if err != nil {
        return
    }
	err = conn.Close()
	conn, _ = net.Dial("udp", "10.96.129.4:60001")
	_, err = conn.Write([]byte(json))
    _ = conn.Close()
	return
}

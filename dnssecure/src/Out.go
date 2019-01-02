package main

import (
	"math/rand"
	"net"
	"time"
)

func  GetRandomString(l int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyz"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

func SendUDP(taskID string, randString string, status string) (err error) {
	conn, err := net.Dial("udp", "10.96.129.4:60001")
	content := taskID + " "  +  randString +" "+ status
	if err != nil {
        return
    }
    _, err = conn.Write([]byte(content))
    if err != nil {
        return
    }
	err = conn.Close()
	return
}

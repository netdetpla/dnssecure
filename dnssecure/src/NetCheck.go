package main

import (
	"fmt"
	"github.com/sparrc/go-ping"
	"time"
)

func PingCheck(addr string) (checkFlag int, err error){
	pinger, err := ping.NewPinger(addr)
	if err != nil {
		fmt.Println(err.Error())
		return 0, err
	}
	pinger.SetPrivileged(true)
	pinger.Count = 1
	pinger.Timeout = time.Duration(2 * time.Second)
	pinger.Run()                 // blocks until finished
	stats := pinger.Statistics() // get send/receive/rtt stats
	if stats.PacketLoss > 0 {
		checkFlag = 0
	} else {
		checkFlag = 1
	}
	return
}

func MultiPingCheck(addr string, count int) (bool, error) {
	recvCount := 0
	var err error
	for i := 0; i < count; i++ {
		var tmpCount int
		tmpCount, err = PingCheck(addr)
		recvCount += tmpCount
	}
	return recvCount != 0, err
}

func NetCheck() (netCheckFlag bool, err error) {
	netCheckFlag, err = MultiPingCheck("8.8.8.8", 10)
	if err != nil {
		return false, err
	}
	if netCheckFlag {
		return
	}
	netCheckFlag, err = MultiPingCheck("114.114.114.114", 10)
	if err != nil {
		return false, err
	}
	return
}

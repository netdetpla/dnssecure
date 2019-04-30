package main

import (
	"fmt"
	"github.com/sparrc/go-ping"
	"time"
)

var endChan = make(chan bool, 1)

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
	fmt.Println(stats)
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

func NetCheckProcess() {
	netCheckFlag, err := MultiPingCheck("8.8.8.8", 10)
	//netCheckFlag, err := MultiPingCheck("2.3.4.2", 60)
	if err != nil {
		endChan <- false
		return
	}
	if netCheckFlag {
		endChan <- true
		return
	}
	netCheckFlag, err = MultiPingCheck("114.114.114.114", 10)
	if err != nil {
		endChan <- false
		return 
	}
	endChan <- netCheckFlag
}

func NetCheck() (netCheckFlag int) {
	go NetCheckProcess()
	clock := time.NewTimer(time.Second)
	for {
		clock.Reset(time.Second * 60)
		select {
		case netFlag := <-endChan:
			fmt.Println("ping check end")
			if netFlag {
				return 1
			} else {
				return 0
			}
		case <-clock.C:
			fmt.Println("ping thread timeout")
			return -1
		}
	}
}

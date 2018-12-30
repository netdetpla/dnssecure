package main

import (
	"github.com/sparrc/go-ping"
	"time"
)

func PingCheck(host string) (packetLoss float64, err error) {
	pinger, err := ping.NewPinger(host)
	if err != nil {
		return 100, err
	}
	pinger.SetPrivileged(true)
	pinger.Count = 10
	pinger.Timeout = time.Duration(1*time.Second)
	pinger.Run()
	stats := pinger.Statistics()
	packetLoss = stats.PacketLoss
	return
}

func NetCheck() (netCheckFlag bool, err error) {
	netCheckFlag = false
	packetLoss, err := PingCheck("8.8.8.8")
	if err != nil {
		return false, err
	}
	if packetLoss <= 90 {
		netCheckFlag = true
		return
	}
	packetLoss, err = PingCheck("114.114.114.114")
	if err != nil {
		return false, err
	}
	if packetLoss <= 90 {
		netCheckFlag = true
		return
	}
	return
}

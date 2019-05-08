package main

import (
	"fmt"
	"github.com/miekg/dns"
	"runtime"
	"strings"
	"time"
)

var quit = make(chan error)
var ctrl = make(chan int, 100)

func ParseRR(rrs []dns.RR) (as []string, cNames []string) {
	fmt.Println(rrs)
	for _, rr := range rrs {
		rrElements := strings.Split(rr.String(), "\t")
		fmt.Println(rrElements)
		if len(rrElements) == 5 {
			if rrElements[3] == "CNAME" {
				cName := string([]rune(rrElements[4])[:len(rrElements[4])-1])
				cNames = append(cNames, cName)
			} else if rrElements[3] == "A" {
				as = append(as, rrElements[4])
			}
		}
	}
	fmt.Println(as)
	fmt.Println(cNames)
	return
}

func SendDNSQuery(record *Record) {
	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn(record.domain), dns.TypeA)
	errCount := 3
Start:
	client := dns.Client{Net: "udp", Timeout: 120 * time.Second}
	in, _, err := client.Exchange(m, record.reServer+":53")
	//in, err := dns.Exchange(m, record.reServer+":53")
	if err != nil {
		fmt.Println(err.Error())
		if errCount == 0 {
			record.timeoutFlag = true
			fmt.Println("time out")
			<-ctrl
			quit <- nil
			return
		} else {
			errCount--
			goto Start
		}
	} else {
		record.timeoutFlag = false
		record.detectAs, record.detectCNames = ParseRR(in.Answer)
	}
	<-ctrl
	quit <- nil
	return
}

func ControlDNSQueryRoutine(tasks *Task) (err error) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	for _, record := range tasks.records {
		ctrl <- 1
		go SendDNSQuery(record)
	}
	for i := 0; i < len(tasks.records); i++ {
		err = <-quit
	}
	close(quit)
	return err
}

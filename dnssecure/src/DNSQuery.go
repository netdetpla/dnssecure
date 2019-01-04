package main

import (
	"fmt"
	"github.com/miekg/dns"
	"runtime"
	"strings"
)

var quit = make(chan error)
var ctrl = make(chan int, 100)

func ParseRR(rrs []dns.RR) (as []string, cNames []string) {
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
	InfoLog("ParseRR")
	fmt.Println(as)
	fmt.Println(cNames)
	return
}

func SendDNSQuery(record *Record) {
	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn(record.rightRecord.domain), dns.TypeA)
	errCount := 3
Start:
	in, err := dns.Exchange(m, record.reServer+":53")
	if err != nil {
		if errCount == 0 {
			record.timeoutFlag = true
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

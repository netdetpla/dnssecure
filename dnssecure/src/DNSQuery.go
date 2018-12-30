package main

import (
	"github.com/miekg/dns"
	"runtime"
	"strings"
	"fmt"
)

var quit = make(chan error)

func ParseRR(rrs []dns.RR) (as []string, cNames []string) {
	for _, rr := range rrs {
		rrElements := strings.Split(rr.String(), "\t")
		fmt.Println(rrElements)
		if len(rrElements) == 5 {
			if rrElements[3] == "CNAME" {
				cName := string([]rune(rrElements[4])[:len(rrElements[4]) - 1])
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
	in, err := dns.Exchange(m, record.reServer + ":53")
	if err != nil {
		if strings.Index(err.Error(), "timeout") >= 0 {
			record.timeoutFlag = true
		} else {
			quit <- err
		}
	} else {
		record.timeoutFlag = false
		record.detectAs, record.detectCNames = ParseRR(in.Answer)
	}
	quit <- nil
}

func ControlDNSQueryRoutine(tasks *Task) (err error){
	runtime.GOMAXPROCS(runtime.NumCPU())
	for _, record := range tasks.records {
		go SendDNSQuery(record)
	}
	for i := 0; i < len(tasks.records); i++ {
		err = <- quit
	}
	close(quit)
	return err
}

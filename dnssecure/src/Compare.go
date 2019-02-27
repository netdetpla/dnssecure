package main

import (
	"fmt"
	"runtime"
)

const INIT = 0
const TRUE = 1
const FALSE = -1

var compareQuit = make(chan int)

func CompareList(detects []string, rights []string) (correctFlag int) {
	InfoLog("Compare")
	fmt.Println(detects)
	fmt.Println(rights)
	correctFlag = TRUE
	if !CheckEmptyStr(detects) {
		correctFlag = FALSE
		return
	}

	for _, detect := range detects {
		singleCorrectFlag := false
		for _, right := range rights {
			fmt.Println(detect, right)
			if detect == right {
				singleCorrectFlag = true
				break
			}
		}
		fmt.Println(singleCorrectFlag)
		if !singleCorrectFlag {
			correctFlag = FALSE
		}
	}
	fmt.Println(correctFlag)
	return
}

func CheckEmptyStr(strList []string) (isEmpty bool){
	isEmpty = false
	for _, s := range strList {
		if len(s) != 0 {
			isEmpty = true
			return
		}
	}
	return
}

func Compare(record *Record) {
	compareAFlag := CheckEmptyStr(record.rightRecord.rightAs)
	compareCNameFlag := CheckEmptyStr(record.rightRecord.rightCNames)
	detectAFlag := CheckEmptyStr(record.detectAs)
	detectCNameFlag := CheckEmptyStr(record.detectCNames)	
	//比对字段类型
	//A/CNAME
	if compareAFlag && compareCNameFlag {
		record.compareType = "A/CNAME"
	} else if compareAFlag {
		record.compareType = "A"
	} else if compareCNameFlag {
		record.compareType = "CNAME"
	}
	//未获取到配置
	if !compareAFlag && !compareCNameFlag {
		record.result = "0-00-1-0-00"
		return
	}
	//查询超时
	if record.timeoutFlag {
		record.result = "0-00-0-0-00"
		return
	}
	//无效应答
	if !detectAFlag && !detectCNameFlag {
		record.result = "0-11-0-0-10"
		return
	}
	correctAFlag := INIT
	correctCNameFlag := INIT
	//A记录与CNAME均需要比较，未探测到A记录
	if compareAFlag && !detectAFlag && compareCNameFlag && detectCNameFlag {
		correctCNameFlag = CompareList(record.detectCNames, record.rightRecord.rightCNames)
		if correctCNameFlag == TRUE {
			//CNAME正确&A记录空
			record.result = "1-01-1-1-001"
			return
		}
	}
	//A记录与CNAME均需要比较，其余情况
	if compareAFlag && compareCNameFlag {
		//比较A记录
		correctAFlag = CompareList(record.detectAs, record.rightRecord.rightAs)
		//比较CNAME
		correctCNameFlag = CompareList(record.detectCNames, record.rightRecord.rightCNames)
	} else if compareAFlag {
		//比较A记录
		correctAFlag = CompareList(record.detectAs, record.rightRecord.rightAs)
	} else if compareCNameFlag {
		//比较CNAME
		correctCNameFlag = CompareList(record.detectCNames, record.rightRecord.rightCNames)
	} else {
		//无效应答（无法比较）
		record.result = "0-11-0-0-10"
		return
	}
	//结果判断
	if correctAFlag + correctCNameFlag > 0 {
		//比对一致
		record.result = "0-11-1-0-00"
		return
	} else if correctAFlag == FALSE && correctCNameFlag == FALSE {
		//A记录与CNAME均错误
		record.result = "0-11-1-1-11"
		return
	} else if correctAFlag == FALSE {
		//A记录错误
		record.result = "0-11-1-1-10"
		return
	} else if correctCNameFlag == FALSE {
		//CNAME错误
		record.result = "0-11-1-1-01"
		return
	} else {
		record.result = "0-11-0-0-10"
		return
	}
}

func CompareBox(record *Record, index int) {
	Compare(record)
	compareQuit <- index
}

func ControlCompareRoutine(tasks *Task) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	InfoLog("ControlCompareRoutine-records")
	fmt.Println(len(tasks.records))
	for index, record := range tasks.records {
		go CompareBox(record, index)
	}
	for i := 0; i < len(tasks.records); i++ {
		<- compareQuit
	}
	close(compareQuit)
}

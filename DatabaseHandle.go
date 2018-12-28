package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"strings"
)

type RightRecord struct {
	domain      string
	rightAs     []string
	rightCNames []string
}

func SplitAs(AsStr string) (As []string) {
	TempAs := strings.Split(AsStr, ",")
	for _, a := range TempAs {
		As = append(As, strings.Split(a, "-")[0])
	}
	return
}

func SplitCNames(CNamesStr string) (CNames []string) {
	CNames = strings.Split(CNamesStr, ",")
	return
}

func getRightValue(domains []string) (rightRecords []*RightRecord, err error) {
	//TODO 备份
	db, err := sql.Open(
		"mysql",
		"root:123456@tcp(192.168.226.11:3306)/cncert_initiative_probe_system")
	if err != nil {
		return nil, err
	}
	//查询数据
	querySQL := "SELECT CNAME, A FROM domain_library WHERE domain_name=?"
	for _, domain := range domains {
		rows, err := db.Query(querySQL, domain)
		if err != nil {
			return nil, err
		}
		rows.Next()
		var (
			rightCNamesStr string
			rightAsStr     string
		)
		err = rows.Scan(&rightCNamesStr, &rightAsStr)
		if err != nil {
			return nil, err
		}
		rightRecord := new(RightRecord{
			domain:      domain,
			rightAs:     SplitAs(rightAsStr),
			rightCNames: SplitCNames(rightCNamesStr),
		})
		rightRecords = append(rightRecords, rightRecord)

		if err != nil {
			os.Exit(1)
		}

	}
	return
}

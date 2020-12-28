package db

import (
	"fmt"
	"github.com/CarosDrean/api-results.git/query"
	"log"
)

func GetResultService(idService string, idExam string, idResult string) string {
	item := ""
	tsql := fmt.Sprintf(query.ResultService["get"].Q, idService, idExam, idResult)
	rows, err := DB.Query(tsql)

	if err != nil {
		fmt.Println("Error reading rows: " + err.Error())
		return item
	}
	for rows.Next() {
		err := rows.Scan(&item)
		if err != nil {
			log.Println(err)
		}
	}
	defer rows.Close()
	return item
}

package db

import (
	"fmt"
	"github.com/CarosDrean/api-results.git/query"
	"log"
)

type ResultDB struct {}

func (db ResultDB) GetService(idService string, idExam string, idResult string) (string, error) {
	item := ""
	tsql := fmt.Sprintf(query.ResultService["get"].Q, idService, idExam, idResult)
	rows, err := DB.Query(tsql)

	if err != nil {
		fmt.Println("Error reading rows: " + err.Error())
		return item, err
	}
	for rows.Next() {
		err := rows.Scan(&item)
		if err != nil {
			log.Println(err)
		}
	}
	defer rows.Close()
	return item, nil
}

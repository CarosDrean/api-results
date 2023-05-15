package db

import (
	"fmt"
	"github.com/CarosDrean/api-results.git/query"
	"log"
)

type StatusGenerateBD struct{}

func (db StatusGenerateBD) GetStatusGenerate(idService string) (string, error) {
	item := ""
	tsql := fmt.Sprintf(query.ResultService["getStatusLiquid"].Q, idService)
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

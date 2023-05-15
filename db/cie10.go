package db

import (
	"database/sql"
	"fmt"
	"github.com/CarosDrean/api-results.git/models"
	"github.com/CarosDrean/api-results.git/query"
)

type CIE10DB struct{}

func (db CIE10DB) GetAll() ([]models.CIE10, error) {
	res := make([]models.CIE10, 0)

	tsql := fmt.Sprintf(query.CIE10["list"].Q)
	rows, err := DB.Query(tsql)

	err = db.scan(rows, err, &res, "CIE10 DB", "GetAll")
	if err != nil {
		return res, err
	}
	defer rows.Close()
	return res, nil
}

func (db CIE10DB) scan(rows *sql.Rows, err error, res *[]models.CIE10, ctx string, situation string) error {
	var item models.CIE10
	if err != nil {
		checkError(err, situation, ctx, "Reading rows")
		return err
	}
	for rows.Next() {
		err := rows.Scan(&item.ID, &item.Description, &item.Description2)
		if err != nil {
			checkError(err, situation, ctx, "Scan rows")
			return err
		} else {
			*res = append(*res, item)
		}
	}
	return nil
}

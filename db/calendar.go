package db

import (
	"database/sql"
	"fmt"
	"github.com/CarosDrean/api-results.git/models"
	"github.com/CarosDrean/api-results.git/query"
)

type CalendarDB struct {}

func (db CalendarDB) GetService(idService string) (models.Calendar, error) {
	res := make([]models.Calendar, 0)

	tsql := fmt.Sprintf(query.Calendar["getServiceID"].Q, idService)
	rows, err := DB.Query(tsql)

	err = db.scan(rows, err, &res, "Calendar DB", "GetAll")
	if err != nil {
		return models.Calendar{}, err
	}
	defer rows.Close()
	return res[0], nil
}

func (db CalendarDB) scan(rows *sql.Rows, err error, res *[]models.Calendar, ctx string, situation string) error {
	var item models.Calendar
	if err != nil {
		checkError(err, situation, ctx, "Reading rows")
		return err
	}
	for rows.Next() {
		err := rows.Scan(&item.ID, &item.ServiceID, &item.CalendarStatusID)
		if err != nil {
			checkError(err, situation, ctx, "Scan rows")
			return err
		} else {
			*res = append(*res, item)
		}
	}
	return nil
}

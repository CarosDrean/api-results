package db

import (
	"database/sql"
	"fmt"
	"github.com/CarosDrean/api-results.git/models"
	"github.com/CarosDrean/api-results.git/query"
)

type LocationDB struct{}

func (db LocationDB) GetAllOrganizationID(id string) ([]models.Location, error) {
	res := make([]models.Location, 0)

	tsql := fmt.Sprintf(query.Location["getOrganizationID"].Q, id)

	rows, err := DB.Query(tsql)

	err = db.scan(rows, err, &res, "Location DB", "GetAll")
	if err != nil {
		return res, err
	}
	defer rows.Close()
	return res, nil
}

func (db LocationDB) Get(id string) (models.Location, error) {
	res := make([]models.Location, 0)

	tsql := fmt.Sprintf(query.Location["get"].Q, id)
	rows, err := DB.Query(tsql)

	err = db.scan(rows, err, &res, "Location DB", "GetAll")
	if err != nil {
		return models.Location{}, err
	}
	defer rows.Close()
	return res[0], err
}

func (db LocationDB) scan(rows *sql.Rows, err error, res *[]models.Location, ctx string, situation string) error {
	var item models.Location
	if err != nil {
		checkError(err, situation, ctx, "Reading rows")
		return err
	}
	for rows.Next() {
		err := rows.Scan(&item.ID, &item.OrganizationID, &item.Name, &item.IsDeleted)
		if err != nil {
			checkError(err, situation, ctx, "Scan rows")
			return err
		} else if item.IsDeleted != 1 {
			*res = append(*res, item)
		}
	}
	return nil
}

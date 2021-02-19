package db

import (
	"database/sql"
	"fmt"
	"github.com/CarosDrean/api-results.git/models"
	"github.com/CarosDrean/api-results.git/query"
)

type ComponentDB struct {}

func (db ComponentDB) GetAllCategoryId(idCategory string) ([]models.Component, error){
	res := make([]models.Component, 0)

	tsql := fmt.Sprintf(query.Component["getCategory"].Q, idCategory)
	rows, err := DB.Query(tsql)

	err = db.scan(rows, err, &res, "Component DB", "GetAll")
	if err != nil {
		return res, err
	}
	defer rows.Close()
	return res, nil
}

func (db ComponentDB) GetComponentProtocolId(idProtocol string) ([]models.Component, error){
	res := make([]models.Component, 0)
	var item models.Component
	tsql := fmt.Sprintf(query.Component["listComponent"].Q, idProtocol)
	rows, err := DB.Query(tsql)
	if err != nil {
		checkError(err, "GetComponentProtocolId", "db", "Reading rows")
		return res, err
	}

	for rows.Next() {
		err = rows.Scan(&item.ID, &item.Name)
		if err != nil {
			checkError(err, "GetComponentProtocolId", "db", "scan rows")
		} else {
			res = append(res, item)
		}
	}
	defer rows.Close()
	return res, err
}



func (db ComponentDB) scan(rows *sql.Rows, err error, res *[]models.Component, ctx string, situation string) error {
	var item models.Component
	if err != nil {
		checkError(err, situation, ctx, "Reading rows")
		return err
	}
	for rows.Next() {
		err := rows.Scan(&item.ID, &item.Name, &item.CategoryID, &item.IsDeleted)
		if err != nil {
			checkError(err, situation, ctx, "Scan rows")
			return err
		} else {
			*res = append(*res, item)
		}
	}
	return nil
}



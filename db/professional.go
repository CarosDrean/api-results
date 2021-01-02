package db

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/CarosDrean/api-results.git/models"
	"github.com/CarosDrean/api-results.git/query"
)

type ProfessionalDB struct {}

func (db ProfessionalDB) GetAll() ([]models.Professional, error) {
	res := make([]models.Professional, 0)

	tsql := fmt.Sprintf(query.Professional["list"].Q)
	rows, err := DB.Query(tsql)

	err = db.scan(rows, err, &res, "Profession DB", "GetAll")
	if err != nil {
		return res, err
	}
	defer rows.Close()
	return res, err
}

func (db ProfessionalDB) Get(id string) (models.Professional, error) {
	res := make([]models.Professional, 0)

	tsql := fmt.Sprintf(query.Professional["get"].Q, id)
	rows, err := DB.Query(tsql)

	err = db.scan(rows, err, &res, "Profession DB", "Get")
	if err != nil {
		return models.Professional{}, err
	}
	if len(res) == 0 {
		return models.Professional{}, nil
	}
	defer rows.Close()
	return res[0], nil
}

func (db ProfessionalDB) Create(item models.Professional) (int64, error) {
	ctx := context.Background()
	tsql := fmt.Sprintf(query.Professional["insert"].Q)

	result, err := DB.ExecContext(
		ctx,
		tsql,
		sql.Named("v_PersonId", item.PersonID),
		sql.Named("i_ProfessionId", item.ProfessionID),
		sql.Named("v_ProfessionalCode", item.Code),
		sql.Named("i_IsDeleted", 0))
	if err != nil {
		return -1, err
	}
	return result.RowsAffected()
}

func (db ProfessionalDB) Update(id string, item models.Professional) (int64, error) {
	ctx := context.Background()
	tsql := fmt.Sprintf(query.Professional["update"].Q)

	result, err := DB.ExecContext(
		ctx,
		tsql,
		sql.Named("ID", id),
		sql.Named("i_ProfessionId", item.ProfessionID),
		sql.Named("v_ProfessionalCode", item.Code),
		sql.Named("i_IsDeleted", 0))
	if err != nil {
		return -1, err
	}
	return result.RowsAffected()
}

func (db ProfessionalDB) Delete(id string) (int64, error) {
	ctx := context.Background()
	tsql := fmt.Sprintf(query.Professional["delete"].Q)
	result, err := DB.ExecContext(
		ctx,
		tsql,
		sql.Named("ID", id))
	if err != nil {
		return -1, err
	}
	return result.RowsAffected()
}

func (db ProfessionalDB) scan(rows *sql.Rows, err error, res *[]models.Professional, ctx string, situation string) error {
	var item models.Professional
	if err != nil {
		checkError(err, situation, ctx, "Reading rows")
		return err
	}
	for rows.Next() {
		err := rows.Scan(&item.PersonID, &item.ProfessionID, &item.Code, &item.IsDeleted)
		if err != nil {
			checkError(err, situation, ctx, "Scan rows")
			return err
		} else {
			*res = append(*res, item)
		}
	}
	return nil
}
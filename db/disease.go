package db

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/CarosDrean/api-results.git/constants"
	"github.com/CarosDrean/api-results.git/models"
	"github.com/CarosDrean/api-results.git/query"
)

type DiseaseDB struct {}

func (db DiseaseDB) GetAll()([]models.Disease, error) {
	res := make([]models.Disease, 0)

	tsql := fmt.Sprintf(query.Disease["list"].Q)
	rows, err := DB.Query(tsql)

	err = db.scan(rows, err, &res, "Disease DB", "GetAll")
	if err != nil {
		return res, err
	}
	defer rows.Close()
	return res, nil
}

func (db DiseaseDB) Create(item models.Disease) (int64, error) {
	ctx := context.Background()
	tsql := fmt.Sprintf(query.Disease["insert"].Q)
	sqdb :=SequentialDB{}
	sequentialID := sqdb.NextSequentialId(constants.IdNode, constants.IdDiseaseTable)
	newId := sqdb.NewID(constants.IdNode, sequentialID, constants.PrefixDisease)
	item.ID = newId
	result, err := DB.ExecContext(
		ctx,
		tsql,
		sql.Named("v_DiseasesId", item.ID),
		sql.Named("v_CIE10Id", item.CIE10ID),
		sql.Named("v_Name", item.Name),
		sql.Named("i_IsDeleted", 0))
	if err != nil {
		return -1, err
	}
	return result.RowsAffected()
}

func (db DiseaseDB) Update(id string, item models.Disease)(int64, error){
	ctx := context.Background()
	tsql := fmt.Sprintf(query.Sequential["update"].Q)
	result, err := DB.ExecContext(
		ctx,
		tsql,
		sql.Named("ID", id),
		sql.Named("v_CIE10Id", item.CIE10ID),
		sql.Named("v_Name", item.Name),
		sql.Named("i_IsDeleted", 0))
	if err != nil {
		return -1, err
	}
	return result.RowsAffected()
}

func (db DiseaseDB) Delete(id string) (int64, error) {
	ctx := context.Background()
	tsql := fmt.Sprintf(query.Disease["delete"].Q)
	result, err := DB.ExecContext(
		ctx,
		tsql,
		sql.Named("ID", id))
	if err != nil {
		fmt.Println(err)
		return -1, err
	}
	return result.RowsAffected()
}

func (db DiseaseDB) scan(rows *sql.Rows, err error, res *[]models.Disease, ctx string, situation string) error {
	var item models.Disease
	if err != nil {
		checkError(err, situation, ctx, "Reading rows")
		return err
	}
	for rows.Next() {
		err := rows.Scan(&item.ID, &item.CIE10ID, &item.Name, &item.IsDeleted)
		if err != nil {
			checkError(err, situation, ctx, "Scan rows")
			return err
		} else {
			*res = append(*res, item)
		}
	}
	return nil
}

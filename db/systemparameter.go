package db

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/CarosDrean/api-results.git/models"
	"github.com/CarosDrean/api-results.git/query"
)

type SystemParameterDB struct {}

func (db SystemParameterDB) GetAllByGroupID(idGroup string) ([]models.SystemParameter, error) {
	res := make([]models.SystemParameter, 0)

	tsql := fmt.Sprintf(query.SystemParameter["getGroup"].Q, idGroup)
	rows, err := DB.Query(tsql)

	err = db.scan(rows, err, &res, "SystemParameter DB", "GetAll")
	if err != nil {
		return res, err
	}
	defer rows.Close()
	return res, nil
}

func (db SystemParameterDB) Create(item models.SystemParameter) (int64, error) {
	ctx := context.Background()
	tsql := fmt.Sprintf(query.SystemParameter["insert"].Q)
	result, err := DB.ExecContext(
		ctx,
		tsql,
		sql.Named("i_GroupId", item.GroupID),
		sql.Named("i_ParameterId", item.ParameterID),
		sql.Named("v_Value1", item.Value1),
		sql.Named("i_IsDeleted", 0))
	if err != nil {
		return -1, err
	}
	return result.RowsAffected()
}

func (db SystemParameterDB) Update(item models.SystemParameter)(int64, error){
	ctx := context.Background()
	tsql := fmt.Sprintf(query.SystemParameter["update"].Q)
	result, err := DB.ExecContext(
		ctx,
		tsql,
		sql.Named("ID", item.GroupID),
		sql.Named("IDT", item.ParameterID),
		sql.Named("v_Value1", item.Value1),
		sql.Named("i_IsDeleted", 0))
	if err != nil {
		return -1, err
	}
	return result.RowsAffected()
}

func (db SystemParameterDB) Delete(item models.SystemParameter) (int64, error) {
	ctx := context.Background()
	tsql := fmt.Sprintf(query.SystemParameter["delete"].Q)
	result, err := DB.ExecContext(
		ctx,
		tsql,
		sql.Named("ID", item.GroupID),
		sql.Named("IDT", item.ParameterID))
	if err != nil {
		fmt.Println(err)
		return -1, err
	}
	return result.RowsAffected()
}

func (db SystemParameterDB) scan(rows *sql.Rows, err error, res *[]models.SystemParameter, ctx string, situation string) error {
	var item models.SystemParameter
	if err != nil {
		checkError(err, situation, ctx, "Reading rows")
		return err
	}
	for rows.Next() {
		err := rows.Scan(&item.GroupID, &item.ParameterID, &item.Value1, &item.IsDeleted)
		if err != nil {
			checkError(err, situation, ctx, "Scan rows")
			return err
		} else {
			*res = append(*res, item)
		}
	}
	return nil
}

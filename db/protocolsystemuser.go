package db

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/CarosDrean/api-results.git/constants"
	"github.com/CarosDrean/api-results.git/models"
	"github.com/CarosDrean/api-results.git/query"
)

type ProtocolSystemUserDB struct {}

func (db ProtocolSystemUserDB) GetAllSystemUserID(id string) ([]models.ProtocolSystemUser, error){
	res := make([]models.ProtocolSystemUser, 0)

	tsql := fmt.Sprintf(query.ProtocolSystemUser["getSystemUserID"].Q, id)
	rows, err := DB.Query(tsql)

	err = db.scan(rows, err, &res, "ProtocolSU DB", "GetAll")
	if err != nil {
		return res, err
	}
	defer rows.Close()
	return res, nil
}

func (db ProtocolSystemUserDB) Get(id string) (models.ProtocolSystemUser, error){
	res := make([]models.ProtocolSystemUser, 0)

	tsql := fmt.Sprintf(query.ProtocolSystemUser["get"].Q, id)
	rows, err := DB.Query(tsql)

	err = db.scan(rows, err, &res, "ProtocolSU DB", "GetAll")
	if err != nil {
		return models.ProtocolSystemUser{}, err
	}
	if len(res) == 0 {
		return models.ProtocolSystemUser{}, nil
	}
	defer rows.Close()
	return res[0], nil
}

func (db ProtocolSystemUserDB) Create(item models.ProtocolSystemUser) (int64, error) {
	ctx := context.Background()
	tsql := fmt.Sprintf(query.ProtocolSystemUser["insert"].Q)

	sqdb := SequentialDB{}
	sequentialID := sqdb.NextSequentialId(constants.IdNode, constants.IdProtocolSystemUserTable)
	newId := sqdb.NewID(constants.IdNode, sequentialID, constants.PrefixProtocolSystemUser)
	item.ID = newId

	result, err := DB.ExecContext(
		ctx,
		tsql,
		sql.Named("v_ProtocolSystemUserId", item.ID),
		sql.Named("i_SystemUserId", item.SystemUserID),
		sql.Named("v_ProtocolId", item.ProtocolID))
	if err != nil {
		return -1, err
	}
	return result.RowsAffected()
}

func (db ProtocolSystemUserDB) scan(rows *sql.Rows, err error, res *[]models.ProtocolSystemUser, ctx string, situation string) error {
	var item models.ProtocolSystemUser
	if err != nil {
		checkError(err, situation, ctx, "Reading rows")
		return err
	}
	for rows.Next() {
		err := rows.Scan(&item.ID, &item.SystemUserID, &item.ProtocolID)
		if err != nil {
			checkError(err, situation, ctx, "Scan rows")
			return err
		} else {
			*res = append(*res, item)
		}
	}
	return nil
}


package db

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/CarosDrean/api-results.git/constants"
	"github.com/CarosDrean/api-results.git/models"
	"github.com/CarosDrean/api-results.git/query"
	"log"
)

func GetProtocolSystemUserWidthSystemUserID(id string) []models.ProtocolSystemUser{
	res := make([]models.ProtocolSystemUser, 0)
	var item models.ProtocolSystemUser

	tsql := fmt.Sprintf(query.ProtocolSystemUser["getSystemUserID"].Q, id)
	rows, err := DB.Query(tsql)

	if err != nil {
		fmt.Println("Error reading rows: " + err.Error())
		return res
	}
	for rows.Next(){
		err := rows.Scan(&item.ID, &item.SystemUserID, &item.ProtocolID)
		if err != nil {
			log.Println(err)
		} else{
			res = append(res, item)
		}
	}
	defer rows.Close()
	return res
}

func GetProtocolSystemUser(id string) []models.ProtocolSystemUser{
	res := make([]models.ProtocolSystemUser, 0)
	var item models.ProtocolSystemUser

	tsql := fmt.Sprintf(query.ProtocolSystemUser["get"].Q, id)
	rows, err := DB.Query(tsql)

	if err != nil {
		fmt.Println("Error reading rows: " + err.Error())
		return res
	}
	for rows.Next(){
		err := rows.Scan(&item.ID, &item.SystemUserID, &item.ProtocolID)
		if err != nil {
			log.Println(err)
		} else{
			res = append(res, item)
		}
	}
	defer rows.Close()
	return res
}

func CreateProtocolSystemUser(item models.ProtocolSystemUser) (int64, error) {
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

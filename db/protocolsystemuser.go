package db

import (
	"fmt"
	"github.com/CarosDrean/api-results.git/models"
	"log"
)

func GetProtocolSystemUserWidthSystemUserID(id string) []models.ProtocolSystemUser{
	res := make([]models.ProtocolSystemUser, 0)
	var item models.ProtocolSystemUser

	tsql := fmt.Sprintf(QueryProtocolSystemUser["getSystemUserID"].Q, id)
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

	tsql := fmt.Sprintf(QueryProtocolSystemUser["get"].Q, id)
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

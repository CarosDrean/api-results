package db

import (
	"fmt"
	"github.com/CarosDrean/api-results.git/models"
	"log"
)

func GetSystemParametersByGroupID(idGroup string) []models.SystemParameter {
	res := make([]models.SystemParameter, 0)
	var item models.SystemParameter

	tsql := fmt.Sprintf(querySystemParameter["getGroup"].Q, idGroup)
	rows, err := DB.Query(tsql)

	if err != nil {
		fmt.Println("Error reading rows: " + err.Error())
		return res
	}
	for rows.Next(){
		err := rows.Scan(&item.GroupID, &item.ParameterID, &item.Value1)
		if err != nil {
			log.Println(err)
		} else{
			res = append(res, item)
		}
	}
	defer rows.Close()
	return res
}

package db

import (
	"fmt"
	"github.com/CarosDrean/api-results.git/models"
	"github.com/CarosDrean/api-results.git/query"
	"log"
)

func GetCalendarService(idService string) models.Calendar {
	res := make([]models.Calendar, 0)
	var item models.Calendar

	tsql := fmt.Sprintf(query.Calendar["getServiceID"].Q, idService)
	rows, err := DB.Query(tsql)

	if err != nil {
		fmt.Println("Error reading rows: " + err.Error())
		return res[0]
	}
	for rows.Next(){
		err := rows.Scan(&item.ID, &item.ServiceID, &item.CalendarStatusID)
		if err != nil {
			log.Println(err)
			return res[0]
		} else{
			res = append(res, item)
		}
	}
	defer rows.Close()
	return item
}

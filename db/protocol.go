package db

import (
	"fmt"
	"github.com/CarosDrean/api-results.git/models"
	"log"
)

func GetProtocolsWidthLocation(id string) models.Protocol {
	res := make([]models.Protocol, 0)
	var item models.Protocol

	tsql := fmt.Sprintf(QueryProtocol["getLocation"].Q, id)
	rows, err := DB.Query(tsql)

	if err != nil {
		fmt.Println("Error reading rows: " + err.Error())
		return res[0]
	}
	for rows.Next(){
		err := rows.Scan(&item.ID, &item.Name, &item.OrganizationID, &item.LocationID)
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

func GetProtocol(id string) models.Protocol {
	res := make([]models.Protocol, 0)
	var item models.Protocol

	tsql := fmt.Sprintf(QueryProtocol["get"].Q, id)
	rows, err := DB.Query(tsql)

	if err != nil {
		fmt.Println("Error reading rows: " + err.Error())
		return res[0]
	}
	for rows.Next(){
		err := rows.Scan(&item.ID, &item.Name, &item.OrganizationID, &item.LocationID)
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

package db

import (
	"fmt"
	"github.com/CarosDrean/api-results.git/models"
	"log"
)

func GetLocationsWidthOrganizationID(id string) []models.Location{
	res := make([]models.Location, 0)
	var item models.Location

	tsql := fmt.Sprintf(queryLocation["getOrganizationID"].Q, id)
	rows, err := DB.Query(tsql)

	if err != nil {
		fmt.Println("Error reading rows: " + err.Error())
		return res
	}
	for rows.Next(){
		err := rows.Scan(&item.ID, &item.OrganizationID, &item.Name, &item.IsDeleted)
		if err != nil {
			log.Println(err)
		} else{
			res = append(res, item)
		}
	}
	defer rows.Close()
	return res
}

func GetLocation(id string) []models.Location{
	res := make([]models.Location, 0)
	var item models.Location

	tsql := fmt.Sprintf(queryLocation["get"].Q, id)
	rows, err := DB.Query(tsql)

	if err != nil {
		fmt.Println("Error reading rows: " + err.Error())
		return res
	}
	for rows.Next(){
		err := rows.Scan(&item.ID, &item.OrganizationID, &item.Name, &item.IsDeleted)
		if err != nil {
			log.Println(err)
		} else{
			res = append(res, item)
		}
	}
	defer rows.Close()
	return res
}


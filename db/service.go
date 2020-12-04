package db

import (
	"fmt"
	"github.com/CarosDrean/api-results.git/models"
	"log"
)

func GetServiceWidthPersonID(id string) []models.Service{
	res := make([]models.Service, 0)
	var item models.Service

	tsql := fmt.Sprintf(QueryService["getPersonID"].Q, id)
	rows, err := DB.Query(tsql)

	if err != nil {
		fmt.Println("Error reading rows: " + err.Error())
		return res
	}
	for rows.Next(){
		err := rows.Scan(&item.ID, &item.PersonID, &item.ProtocolID, &item.ServiceDate, &item.ServiceStatusId, &item.IsDeleted)
		if err != nil {
			log.Println(err)
		} else{
			res = append(res, item)
		}
	}
	defer rows.Close()
	return res
}

func GetServicesWidthProtocol(id string) []models.Service{
	res := make([]models.Service, 0)
	var item models.Service

	tsql := fmt.Sprintf(QueryService["getProtocol"].Q, id)
	rows, err := DB.Query(tsql)

	if err != nil {
		fmt.Println("Error reading rows: " + err.Error())
		return res
	}
	for rows.Next(){
		err := rows.Scan(&item.ID, &item.PersonID, &item.ProtocolID, &item.ServiceDate, &item.ServiceStatusId, &item.IsDeleted)
		if err != nil {
			log.Println(err)
		} else{
			res = append(res, item)
		}
	}
	defer rows.Close()
	return res
}

func GetServicesWidthProtocolFilter(filter models.Filter) []models.Service{
	res := make([]models.Service, 0)
	var item models.Service

	tsql := fmt.Sprintf(QueryService["getProtocolFilter"].Q, filter.ID, filter.DateFrom, filter.DateTo)
	rows, err := DB.Query(tsql)

	if err != nil {
		fmt.Println("Error reading rows: " + err.Error())
		return res
	}
	for rows.Next(){
		err := rows.Scan(&item.ID, &item.PersonID, &item.ProtocolID, &item.ServiceDate, &item.ServiceStatusId, &item.IsDeleted)
		if err != nil {
			log.Println(err)
		} else{
			res = append(res, item)
		}
	}
	defer rows.Close()
	return res
}

func GetService(id string) []models.Service{
	res := make([]models.Service, 0)
	var item models.Service

	tsql := fmt.Sprintf(QueryService["get"].Q, id)
	rows, err := DB.Query(tsql)

	if err != nil {
		fmt.Println("Error reading rows: " + err.Error())
		return res
	}
	for rows.Next(){
		err := rows.Scan(&item.ID, &item.PersonID, &item.ProtocolID, &item.ServiceDate, &item.ServiceStatusId, &item.IsDeleted)
		if err != nil {
			log.Println(err)
		} else{
			res = append(res, item)
		}
	}
	defer rows.Close()
	return res
}

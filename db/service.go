package db

import (
	"fmt"
	"github.com/CarosDrean/api-results.git/models"
	"log"
)

const (
	NQGetServicePersonID       nameQuery = "getPersonID"
	NQGetServiceProtocol       nameQuery = "getProtocol"
	NQGetServiceProtocolFilter nameQuery = "getProtocolFilter"
	NQGetService               nameQuery = "get"
)

func GetServicesWidthProtocolFilter(filter models.Filter) []models.Service {
	res := make([]models.Service, 0)
	var item models.Service

	tsql := fmt.Sprintf(QueryService[NQGetServiceProtocolFilter].Q, filter.ID, filter.DateFrom, filter.DateTo)
	rows, err := DB.Query(tsql)

	if err != nil {
		fmt.Println("Error reading rows: " + err.Error())
		return res
	}
	for rows.Next() {
		err := rows.Scan(&item.ID, &item.PersonID, &item.ProtocolID, &item.ServiceDate, &item.ServiceStatusId,
			&item.IsDeleted, &item.AptitudeStatusId)
		if err != nil {
			log.Println(err)
		} else {
			res = append(res, item)
		}
	}
	defer rows.Close()
	return res
}

func GetService(id string, nameQuery nameQuery) []models.Service {
	res := make([]models.Service, 0)
	var item models.Service

	tsql := fmt.Sprintf(QueryService[nameQuery].Q, id)
	rows, err := DB.Query(tsql)

	if err != nil {
		fmt.Println("Error reading rows: " + err.Error())
		return res
	}
	for rows.Next() {
		err := rows.Scan(&item.ID, &item.PersonID, &item.ProtocolID, &item.ServiceDate, &item.ServiceStatusId,
			&item.IsDeleted, &item.AptitudeStatusId)
		if err != nil {
			log.Println(err)
		} else if item.IsDeleted != 1 && item.ServiceStatusId == 3 { // verificar servicios no eliminados y culminados
			res = append(res, item)
		}
	}
	defer rows.Close()
	return res
}

// esta funcion es para la peticion de estadisticas
func GetServicesFilter(id string, filter models.Filter) []models.Service {
	res := make([]models.Service, 0)
	var item models.Service

	tsql := fmt.Sprintf(QueryService[NQGetServiceProtocolFilter].Q, id, filter.DateFrom, filter.DateTo)
	rows, err := DB.Query(tsql)

	if err != nil {
		fmt.Println("Error reading rows: " + err.Error())
		return res
	}
	for rows.Next() {
		err := rows.Scan(&item.ID, &item.PersonID, &item.ProtocolID, &item.ServiceDate, &item.ServiceStatusId,
			&item.IsDeleted, &item.AptitudeStatusId)
		if err != nil {
			log.Println(err)
		} else if item.IsDeleted != 1 && item.ServiceStatusId == 3 { // verificar servicios no eliminados y culminados
			res = append(res, item)
		}
	}
	defer rows.Close()
	return res
}

package db

import (
	"database/sql"
	"fmt"
	"github.com/CarosDrean/api-results.git/models"
	"github.com/CarosDrean/api-results.git/query"
	"log"
)

const (
	NQGetServicePersonID       = "getPersonID"
	NQGetServiceProtocol       = "getProtocol"
	NQGetServiceProtocolFilter = "getProtocolFilter"
	NQGetService               = "get"
)

func GetServicesFilterDate(filter models.Filter) []models.ServicePatientDiseases {
	res := make([]models.ServicePatientDiseases, 0)
	var service models.Service
	var person models.Person
	var protocol models.Protocol

	tsql := fmt.Sprintf(query.Service["listDiseaseFilter"].Q, filter.DateFrom, filter.DateTo)
	rows, err := DB.Query(tsql)

	if err != nil {
		fmt.Println("Error reading rows: " + err.Error())
		return res
	}
	for rows.Next() {
		var pass sql.NullString
		var birth sql.NullString
		var disease sql.NullString
		var diseaseString string
		err := rows.Scan(&service.ID, &service.PersonID, &service.ProtocolID, &service.ServiceDate, &service.ServiceStatusId,
			&service.IsDeleted, &service.AptitudeStatusId,
			&person.ID, &person.DNI, &pass, &person.Name, &person.FirstLastName, &person.SecondLastName, &person.Mail,
			&person.Sex, &birth, &person.IsDeleted,
			&protocol.ID, &protocol.Name, &protocol.OrganizationID, &protocol.LocationID, &protocol.IsDeleted, &protocol.EsoType,
			&disease)
		if pass.Valid {
			person.Password = pass.String
		} else {
			person.Password = ""
		}
		if birth.Valid {
			person.Birthday = birth.String
		} else {
			person.Birthday = ""
		}
		if disease.Valid {
			diseaseString = disease.String
		} else {
			diseaseString = ""
		}
		if err != nil {
			log.Println(err)
		} else {
			item := models.ServicePatientDiseases{
				ID:               service.ID,
				ServiceDate:      service.ServiceDate,
				PersonID:         service.PersonID,
				ProtocolID:       service.ProtocolID,
				OrganizationID:   protocol.OrganizationID,
				AptitudeStatusId: service.AptitudeStatusId,
				DNI:              person.DNI,
				Name:             person.Name,
				FirstLastName:    person.FirstLastName,
				SecondLastName:   person.SecondLastName,
				Mail:             person.Mail,
				Sex:              person.Sex,
				Birthday:         person.Birthday,
				Disease:          diseaseString,
				EsoType:          protocol.EsoType,
			}
			res = append(res, item)
		}
	}
	defer rows.Close()
	return res
}

func GetServicesWidthProtocolFilter(filter models.Filter) []models.Service {
	res := make([]models.Service, 0)
	var item models.Service

	fmt.Println(filter)

	tsql := fmt.Sprintf(query.Service[NQGetServiceProtocolFilter].Q, filter.ID, filter.DateFrom, filter.DateTo)
	rows, err := DB.Query(tsql)
	fmt.Println(tsql)

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

func GetService(id string, nameQuery string) []models.Service {
	res := make([]models.Service, 0)
	var item models.Service

	tsql := fmt.Sprintf(query.Service[nameQuery].Q, id)
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

	tsql := fmt.Sprintf(query.Service[NQGetServiceProtocolFilter].Q, id, filter.DateFrom, filter.DateTo)
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

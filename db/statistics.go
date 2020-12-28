package db

import (
	"fmt"
	"github.com/CarosDrean/api-results.git/models"
	"github.com/CarosDrean/api-results.git/query"
	"log"
)

func GetStatisticsServiceDiseaseByProtocol(filter models.Filter)[]models.ServicePatientDiseases {
	res := make([]models.ServicePatientDiseases, 0)
	var item models.ServicePatientDiseases

	tsql := fmt.Sprintf(query.Statistics["getDisease"].Q, filter.ID, filter.DateFrom, filter.DateTo)
	rows, err := DB.Query(tsql)

	if err != nil {
		fmt.Println("Error reading rows: " + err.Error())
		return res
	}
	for rows.Next() {
		err := rows.Scan(&item.ID, &item.ServiceDate, &item.PersonID, &item.ProtocolID, &item.AptitudeStatusId,
			&item.DNI, &item.Name, &item.FirstLastName, &item.SecondLastName, &item.Mail, &item.Sex, &item.Birthday,
			&item.Disease, &item.Component, &item.Consulting)
		if err != nil {
			log.Println(err)
		} else {
			res = append(res, item)
		}
	}
	defer rows.Close()
	return res
}

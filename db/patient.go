package db

import (
	"fmt"
	"github.com/CarosDrean/api-results.git/models"
	"log"
)

func GetPatientFromDNI(dni string) []models.Patient {
	res := make([]models.Patient, 0)
	var item models.Patient

	tsql := fmt.Sprintf(QueryPatient["getDNI"].Q, dni)
	rows, err := DB.Query(tsql)

	if err != nil {
		fmt.Println("Error reading rows: " + err.Error())
		return res
	}
	for rows.Next(){
		err := rows.Scan(&item.ID, &item.DNI, &item.Password)
		if err != nil {
			log.Println(err)
			return res
		} else{
			res = append(res, item)
		}
	}
	defer rows.Close()

	return res
}

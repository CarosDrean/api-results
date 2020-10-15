package db

import (
	"database/sql"
	"github.com/CarosDrean/api-results.git/models"
	"log"
)

func GetPatient(dni int) []models.Patient {
	res := make([]models.Patient, 0)
	var item models.Patient
	get := PrepStmtsPatient["get"].Stmt
	err := get.QueryRow(dni).Scan(&item.ID, &item.Name, &item.DNI)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Printf("user: error getting post. Id: %d, err: %v\n", dni, err)
		}
	} else {
		res = append(res, item)
	}
	return res
}

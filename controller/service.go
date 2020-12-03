package controller

import (
	"encoding/json"
	"github.com/CarosDrean/api-results.git/db"
	"github.com/CarosDrean/api-results.git/models"
	"github.com/gorilla/mux"
	"net/http"
)

func GetServicesPatientsWithProtocol(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var params = mux.Vars(r)
	id, _ := params["idProtocol"]

	res := make([]models.ServicePatient, 0)

	services := db.GetServicesWidthProtocol(id)
	for _, e := range services {
		patient := db.GetPatient(e.PersonID)[0]
		item := models.ServicePatient{
			ID:             e.ID,
			PersonID:       patient.ID,
			ServiceDate:    e.ServiceDate,
			DNI:            patient.DNI,
			Name:           patient.Name,
			FirstLastName:  patient.FirstLastName,
			SecondLastName: patient.SecondLastName,
			Mail:           patient.Mail,
		}
		res = append(res, item)
	}

	_ = json.NewEncoder(w).Encode(res)
}

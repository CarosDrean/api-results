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

	services := db.GetService(id, db.NQGetServiceProtocol)
	for _, e := range services {
		patient := db.GetPatient(e.PersonID)[0]
		item := models.ServicePatient{
			ID:             e.ID,
			PersonID:       patient.ID,
			ServiceDate:    e.ServiceDate,
			DNI:            patient.DNI,
			Name:           patient.Name,
			ProtocolID:     e.ProtocolID,
			FirstLastName:  patient.FirstLastName,
			SecondLastName: patient.SecondLastName,
			Mail:           patient.Mail,
			Sex:            patient.Sex,
		}
		res = append(res, item)
	}

	_ = json.NewEncoder(w).Encode(res)
}

func GetServicesPatientsWithOrganization(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var params = mux.Vars(r)
	id, _ := params["id"]

	res := make([]models.ServicePatient, 0)

	// deacuerdo al id de la empresa obtener todos sus protocolos e ir armando el objeto

	protocols := db.GetProtocolsWidthOrganization(id)
	for _, e := range protocols {
		services := db.GetService(e.ID, db.NQGetServiceProtocol)
		for _, s := range services {
			patient := db.GetPatient(s.PersonID)[0]
			item := models.ServicePatient{
				ID:               s.ID,
				ServiceDate:      s.ServiceDate,
				PersonID:         patient.ID,
				ProtocolID:       s.ProtocolID,
				AptitudeStatusId: s.AptitudeStatusId,
				DNI:              patient.DNI,
				Name:             patient.Name,
				FirstLastName:    patient.FirstLastName,
				SecondLastName:   patient.SecondLastName,
				Mail:             patient.Mail,
				Sex:              patient.Sex,
				Birthday:         patient.Birthday,
			}
			res = append(res, item)
		}
	}

	_ = json.NewEncoder(w).Encode(res)
}

func GetServicesPatientsWithProtocolFilter(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var item models.Filter
	_ = json.NewDecoder(r.Body).Decode(&item)

	res := make([]models.ServicePatient, 0)

	services := db.GetServicesWidthProtocolFilter(item)
	for _, e := range services {
		patient := db.GetPatient(e.PersonID)[0]
		item := models.ServicePatient{
			ID:             e.ID,
			PersonID:       patient.ID,
			ServiceDate:    e.ServiceDate,
			ProtocolID:     e.ProtocolID,
			DNI:            patient.DNI,
			Name:           patient.Name,
			FirstLastName:  patient.FirstLastName,
			SecondLastName: patient.SecondLastName,
			Mail:           patient.Mail,
			Sex:            patient.Sex,
		}
		res = append(res, item)
	}

	_ = json.NewEncoder(w).Encode(res)
}

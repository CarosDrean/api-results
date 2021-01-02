package controller

import (
	"encoding/json"
	"github.com/CarosDrean/api-results.git/constants"
	"github.com/CarosDrean/api-results.git/db"
	"github.com/CarosDrean/api-results.git/models"
	"github.com/gorilla/mux"
	"net/http"
)

type ServiceController struct {
	DB db.ServiceDB
}

func (c ServiceController) GetAllDiseaseFilterDate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var item models.Filter
	_ = json.NewDecoder(r.Body).Decode(&item)

	services := c.DB.GetAllDiseaseFilterDate(item)

	_ = json.NewEncoder(w).Encode(services)
}

func (c ServiceController) GetAllPatientsWithProtocol(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var params = mux.Vars(r)
	id, _ := params["idProtocol"]

	res := make([]models.ServicePatient, 0)

	services, _ := c.DB.GetAllProtocol(id)
	for _, e := range services {
		patient, _ := db.PersonDB{}.Get(e.PersonID)
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

func (c ServiceController) GetAllPatientsWithOrganization(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var params = mux.Vars(r)
	id, _ := params["id"]

	res := make([]models.ServicePatient, 0)

	// deacuerdo al id de la empresa obtener todos sus protocolos e ir armando el objeto

	protocols := db.ProtocolDB{}.GetAllOrganization(id)
	for _, e := range protocols {
		services, _ := c.DB.GetAllProtocol(e.ID)
		for _, s := range services {
			patient, _ := db.PersonDB{}.Get(s.PersonID)
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
				Result:           db.GetResultService(s.ID, constants.IdPruebaRapida, constants.IdResultPruebaRapida),
			}
			res = append(res, item)
		}
	}

	_ = json.NewEncoder(w).Encode(res)
}

func (c ServiceController) GetAllPatientsWithOrganizationFilter(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	var item models.Filter
	_ = json.NewDecoder(r.Body).Decode(&item)

	res := make([]models.ServicePatient, 0)

	// deacuerdo al id de la empresa obtener todos sus protocolos e ir armando el objeto

	protocols := db.ProtocolDB{}.GetAllOrganization(item.ID)
	for _, e := range protocols {
		services, _ := c.DB.GetAllProtocolFilter(e.ID, item)
		for _, s := range services {
			patient, _ := db.PersonDB{}.Get(s.PersonID)
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
				Result:           db.GetResultService(s.ID, constants.IdPruebaRapida, constants.IdResultPruebaRapida),
			}
			res = append(res, item)
		}
	}

	_ = json.NewEncoder(w).Encode(res)
}

func (c ServiceController) GetAllPatientsWithProtocolFilter(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var item models.Filter
	_ = json.NewDecoder(r.Body).Decode(&item)

	res := make([]models.ServicePatient, 0)

	services, _ := c.DB.GetAllProtocolFilter(item.ID, item)
	for _, e := range services {
		patient, _ := db.PersonDB{}.Get(e.PersonID)
		item := models.ServicePatient{
			ID:               e.ID,
			ServiceDate:      e.ServiceDate,
			PersonID:         patient.ID,
			ProtocolID:       e.ProtocolID,
			AptitudeStatusId: e.AptitudeStatusId,
			DNI:              patient.DNI,
			Name:             patient.Name,
			FirstLastName:    patient.FirstLastName,
			SecondLastName:   patient.SecondLastName,
			Mail:             patient.Mail,
			Sex:              patient.Sex,
			Result:           db.GetResultService(e.ID, constants.IdPruebaRapida, constants.IdResultPruebaRapida),
		}
		res = append(res, item)
	}

	_ = json.NewEncoder(w).Encode(res)
}

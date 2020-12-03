package controller

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/CarosDrean/api-results.git/db"
	"github.com/CarosDrean/api-results.git/models"
	"github.com/gorilla/mux"
)

func GetPatientsWithProtocol(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var params = mux.Vars(r)
	id, _ := params["idProtocol"]

	res := make([]models.Patient, 0)
	var item models.Patient

	services := db.GetServicesWidthProtocol(id)
	for _, e := range services {
		item = db.GetPatient(e.PersonID)[0]
		res = append(res, item)
	}

	_ = json.NewEncoder(w).Encode(res)
}

func GetPatient(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var params = mux.Vars(r)
	id, _ := params["id"]

	items := db.GetPatient(id)

	_ = json.NewEncoder(w).Encode(items[0])
}

func UpdatePasswordPatient(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var params = mux.Vars(r)
	id, _ := params["id"]
	var patient models.Patient
	_ = json.NewDecoder(r.Body).Decode(&patient)
	_, err := db.UpdatePasswordPatient(id, patient.Password)
	if err != nil {
		log.Println(err)
		return
	}

	_ = json.NewEncoder(w).Encode(patient)
}

func GetPatientFromDNI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var params = mux.Vars(r)
	dni, _ := params["dni"]

	items := db.GetPatientFromDNI(dni)

	_ = json.NewEncoder(w).Encode(items)
}

func GetPatientFromLogin(user models.UserLogin) []models.Patient {
	items := db.GetPatientFromDNI(user.User)
	return items
}

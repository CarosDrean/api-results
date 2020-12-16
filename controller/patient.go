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

	res := make([]models.Person, 0)
	var item models.Person

	services := db.GetService(id, db.NQGetServiceProtocol)
	for _, e := range services {
		item = db.GetPerson(e.PersonID)[0]
		res = append(res, item)
	}

	_ = json.NewEncoder(w).Encode(res)
}

func GetPatient(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var params = mux.Vars(r)
	id, _ := params["id"]

	items := db.GetPerson(id)

	_ = json.NewEncoder(w).Encode(items[0])
}

func UpdatePasswordPatient(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var params = mux.Vars(r)
	id, _ := params["id"]
	var patient models.Person
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

	item := db.GetPersonFromDNI(dni)

	_ = json.NewEncoder(w).Encode(item)
}

func GetPatientFromLogin(user models.UserLogin) []models.Person {
	items := db.GetPersonFromDNI(user.User)
	return items
}

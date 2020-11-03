package controller

import (
	"encoding/json"
	"net/http"

	"github.com/CarosDrean/api-results.git/db"
	"github.com/CarosDrean/api-results.git/models"
	"github.com/gorilla/mux"
)

func GetPatient(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var params = mux.Vars(r)
	id, _ := params["id"]

	items := db.GetPatient(id)

	_ = json.NewEncoder(w).Encode(items[0])
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

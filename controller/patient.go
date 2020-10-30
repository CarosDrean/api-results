package controller

import (
	"encoding/json"
	"net/http"

	"github.com/CarosDrean/api-results.git/db"
	"github.com/CarosDrean/api-results.git/models"
	"github.com/gorilla/mux"
)

func ValidatePatient() {

}

func GetPatientFromDNI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var params = mux.Vars(r)
	dni, _ := params["dni"]

	items := db.GetPatientFromDNI(dni)

	_ = json.NewEncoder(w).Encode(items)
}

func GetPatientFromLogin(user models.UserLogin) []models.Patient {
	items := db.GetPatientFromDNI(user.Email)
	return items
}

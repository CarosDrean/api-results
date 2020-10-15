package controller

import (
	"encoding/json"
	"github.com/CarosDrean/api-results.git/db"
	"github.com/CarosDrean/api-results.git/models"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func InitP(){
	var err error
	for _, sc := range db.PrepStmtsPatient{
		sc.Stmt, err = db.DB.Prepare(sc.Q)
		if err != nil {
			log.Panic(err)
		}
	}
}

func ValidatePatient(){

}

func GetPatientFromDNI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var params = mux.Vars(r)
	dni, _ := params["dni"]

	InitP()
	items := db.GetPatientFromDNI(dni)

	_ = json.NewEncoder(w).Encode(items)
}

func GetPatientFromLogin(user models.UserLogin) []models.Patient {
	InitP()
	items := db.GetPatientFromDNI(user.Email)
	return items
}

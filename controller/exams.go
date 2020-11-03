package controller

import (
	"encoding/json"
	"github.com/CarosDrean/api-results.git/db"
	"github.com/CarosDrean/api-results.git/models"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func GetExams(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	var params = mux.Vars(r)
	id, _ := params["id"]
	res := make([]models.Result, 0)
	var item models.Result

	patients := db.GetPatient(id)
	services := db.GetServiceWidthPersonID(patients[0].ID)

	for i, e := range services {
		item.ID = strconv.Itoa(i)
		item.ServiceDate = e.ServiceDate
		item.ProtocolName = db.GetProtocol(e.ProtocolID).Name
		res = append(res, item)
	}
	_ = json.NewEncoder(w).Encode(res)
}

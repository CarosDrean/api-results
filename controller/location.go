package controller

import (
	"encoding/json"
	"github.com/CarosDrean/api-results.git/db"
	"github.com/gorilla/mux"
	"net/http"
)

func GetLocationsWidthOrganizationID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var params = mux.Vars(r)
	id, _ := params["idOrganization"]

	res := db.GetLocationsWidthOrganizationID(id)

	_ = json.NewEncoder(w).Encode(res)
}

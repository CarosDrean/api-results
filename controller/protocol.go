package controller

import (
	"encoding/json"
	"github.com/CarosDrean/api-results.git/db"
	"github.com/gorilla/mux"
	"net/http"
)

func GetProtocolsWidthLocation(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	var params = mux.Vars(r)
	id, _ := params["idLocation"]

	protocols := db.GetProtocolsWidthLocation(id)

	_ = json.NewEncoder(w).Encode(protocols)
}

func GetProtocolsWidthOrganization(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	var params = mux.Vars(r)
	id, _ := params["idOrganization"]

	protocols := db.GetProtocolsWidthOrganization(id)

	_ = json.NewEncoder(w).Encode(protocols)
}

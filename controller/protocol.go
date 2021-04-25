package controller

import (
	"encoding/json"
	"github.com/CarosDrean/api-results.git/db"
	"github.com/gorilla/mux"
	"net/http"
)

type ProtocolController struct {
	DB db.ProtocolDB
}

func (c ProtocolController) Get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var params = mux.Vars(r)
	id, _ := params["id"]
	protocol, err := c.DB.Get(id)
	if err != nil {
		returnErr(w, err, "obtener todos protocol")
		return
	}
	_ = json.NewEncoder(w).Encode(protocol)
}

func (c ProtocolController) GetAllLocation(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	var params = mux.Vars(r)
	id, _ := params["idLocation"]

	protocols, err := c.DB.GetAllLocation(id)
	if err != nil {
		returnErr(w, err, "obtener location")
		return
	}

	_ = json.NewEncoder(w).Encode(protocols)
}

func (c ProtocolController) GetAllOrganization(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	var params = mux.Vars(r)
	id, _ := params["idOrganization"]


	protocols, err := c.DB.GetAllOrganization(id)
	if err != nil {
		returnErr(w, err, "obtener todos organization")
		return
	}

	_ = json.NewEncoder(w).Encode(protocols)
}

func (c ProtocolController) GetAllOrganizationEmployer(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	var params = mux.Vars(r)
	id, _ := params["idOrganization"]

	protocols, err := c.DB.GetAllOrganizationEmployer(id)
	if err != nil {
		returnErr(w, err, "obtener todos organization")
		return
	}

	_ = json.NewEncoder(w).Encode(protocols)
}


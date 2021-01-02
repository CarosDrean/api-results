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
	protocol := c.DB.Get(id)
	_ = json.NewEncoder(w).Encode(protocol)
}

func (c ProtocolController) GetAllLocation(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	var params = mux.Vars(r)
	id, _ := params["idLocation"]

	protocols := c.DB.GetAllLocation(id)

	_ = json.NewEncoder(w).Encode(protocols)
}

func (c ProtocolController) GetAllOrganization(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	var params = mux.Vars(r)
	id, _ := params["idOrganization"]

	protocols := c.DB.GetAllOrganization(id)

	_ = json.NewEncoder(w).Encode(protocols)
}

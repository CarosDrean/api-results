package controller

import (
	"encoding/json"
	"net/http"

	"github.com/CarosDrean/api-results.git/db"
	"github.com/CarosDrean/api-results.git/models"
	"github.com/gorilla/mux"
)

type PersonController struct {
	DB db.PersonDB
}

func (c PersonController) GetAllProtocol(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var params = mux.Vars(r)
	id, _ := params["idProtocol"]

	res := make([]models.Person, 0)
	var item models.Person

	services, err := db.ServiceDB{}.GetAllProtocol(id)
	if err != nil {
		returnErr(w, err, "obtener todos")
		return
	}
	for _, e := range services {
		item, _ = c.DB.Get(e.PersonID)
		res = append(res, item)
	}

	_ = json.NewEncoder(w).Encode(res)
}

func (c PersonController) Get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var params = mux.Vars(r)
	id, _ := params["id"]

	item, err := c.DB.Get(id)
	if err != nil {
		returnErr(w, err, "obtener")
		return
	}

	_ = json.NewEncoder(w).Encode(item)
}

func (c PersonController) UpdatePassword(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var params = mux.Vars(r)
	id, _ := params["id"]
	var patient models.Person
	_ = json.NewDecoder(r.Body).Decode(&patient)
	_, err := c.DB.UpdatePassword(id, patient.Password)
	if err != nil {
		returnErr(w, err, "updated password")
		return
	}

	_ = json.NewEncoder(w).Encode(patient)
}

func (c PersonController) GetFromDNI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var params = mux.Vars(r)
	dni, _ := params["dni"]

	item, err := c.DB.GetFromDNI(dni)
	if err != nil {
		returnErr(w, err, "obtener por dnni")
		return
	}

	_ = json.NewEncoder(w).Encode(item)
}

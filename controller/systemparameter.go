package controller

import (
	"encoding/json"
	"github.com/CarosDrean/api-results.git/constants"
	"github.com/CarosDrean/api-results.git/db"
	"github.com/CarosDrean/api-results.git/models"
	"net/http"
)

type SystemParameterController struct {
	DB db.SystemParameterDB
}

func (c SystemParameterController) GetConsultingS(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	res, err := c.DB.GetAllByGroupID(constants.IdConsultings)
	if err != nil {
		returnErr(w, err, "obtener todos")
		return
	}

	_ = json.NewEncoder(w).Encode(res)
}

func (c SystemParameterController) Create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var item models.SystemParameter
	_ = json.NewDecoder(r.Body).Decode(&item)
	result, err := c.DB.Create(item)
	if err != nil {
		returnErr(w, err, "crear")
		return
	}

	_ = json.NewEncoder(w).Encode(result)
}

func (c SystemParameterController) Update(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var item models.SystemParameter
	_ = json.NewDecoder(r.Body).Decode(&item)
	result, err := c.DB.Update(item)
	if err != nil {
		returnErr(w, err, "update")
		return
	}

	_ = json.NewEncoder(w).Encode(result)
}

func (c SystemParameterController) Delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var item models.SystemParameter
	_ = json.NewDecoder(r.Body).Decode(&item)
	result, err := c.DB.Delete(item)
	if err != nil {
		returnErr(w, err, "delete")
		return
	}

	_ = json.NewEncoder(w).Encode(result)
}

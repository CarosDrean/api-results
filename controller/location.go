package controller

import (
	"encoding/json"
	"github.com/CarosDrean/api-results.git/db"
	"github.com/gorilla/mux"
	"net/http"
)

type LocationController struct {
	DB db.LocationDB
}

func (c LocationController) GetAllOrganizationID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var params = mux.Vars(r)
	id, _ := params["idOrganization"]

	res, err := c.DB.GetAllOrganizationID(id)
	if err != nil {
		returnErr(w, err, "obtener todos Organization")
		return
	}

	_ = json.NewEncoder(w).Encode(res)
}

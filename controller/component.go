package controller

import (
	"encoding/json"
	"github.com/CarosDrean/api-results.git/db"
	"github.com/gorilla/mux"
	"net/http"
)

type ComponentController struct {
	DB db.ComponentDB
}

func (c ComponentController) GetAllCategoryId(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var params = mux.Vars(r)
	id, _ := params["id"]

	res, err := c.DB.GetAllCategoryId(id)
	if err != nil {
		returnErr(w, err, "obtener todos category")
		return
	}

	_ = json.NewEncoder(w).Encode(res)
}

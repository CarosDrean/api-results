package controller

import (
	"encoding/json"
	"github.com/CarosDrean/api-results.git/db"
	"net/http"
)

type CIE10Controller struct {
	DB db.CIE10DB
}

func (c CIE10Controller) GetAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	res, err := c.DB.GetAll()
	if err != nil {
		returnErr(w, err, "obtener todos")
		return
	}

	_ = json.NewEncoder(w).Encode(res)
}

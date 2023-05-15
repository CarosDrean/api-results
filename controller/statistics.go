package controller

import (
	"encoding/json"
	"github.com/CarosDrean/api-results.git/db"
	"github.com/CarosDrean/api-results.git/models"
	"net/http"
)

type StatisticController struct{}

func (c StatisticController) GetServiceDiseaseByProtocol(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var item models.Filter
	_ = json.NewDecoder(r.Body).Decode(&item)

	items, err := db.StatisticDB{}.GetServiceDiseaseByProtocol(item)
	if err != nil {
		returnErr(w, err, "obtener statistic")
		return
	}

	_ = json.NewEncoder(w).Encode(items)
}

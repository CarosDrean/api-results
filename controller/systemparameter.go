package controller

import (
	"encoding/json"
	"github.com/CarosDrean/api-results.git/constants"
	"github.com/CarosDrean/api-results.git/db"
	"net/http"
)

type SystemParameterController struct {}

func (c SystemParameterController) GetConsultingS(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	res := db.GetSystemParametersByGroupID(constants.IdConsultings)

	_ = json.NewEncoder(w).Encode(res)
}

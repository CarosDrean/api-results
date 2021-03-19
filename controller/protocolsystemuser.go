package controller

import (
	"encoding/json"
	"github.com/CarosDrean/api-results.git/db"
	"github.com/CarosDrean/api-results.git/models"
	"github.com/gorilla/mux"
	"net/http"
)

type ProtocolSystemUserController struct {
	DB db.ProtocolSystemUserDB
}

func (c ProtocolSystemUserController) GetSystemUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var params = mux.Vars(r)
	id, _ := params["idSystemUser"]

	res := make([]models.Protocol, 0)
	var item models.Protocol

	protocolsSystemUser, err := c.DB.GetAllSystemUserID(id)
	if err != nil {
		returnErr(w, err, "obtener todos")
		return
	}
	for _, e := range protocolsSystemUser {
		item, _ = db.ProtocolDB{}.Get(e.ProtocolID)
		res = append(res, item)
	}


	_ = json.NewEncoder(w).Encode(res)
}


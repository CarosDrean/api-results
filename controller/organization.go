package controller

import (
	"encoding/json"
	"github.com/CarosDrean/api-results.git/constants"
	"github.com/CarosDrean/api-results.git/db"
	mid "github.com/CarosDrean/api-results.git/middleware"
	"github.com/CarosDrean/api-results.git/models"
	"github.com/CarosDrean/api-results.git/utils"
	"github.com/gorilla/mux"
	"net/http"
)

func GetOrganizations(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	items := db.GetOrganizations()

	_ = json.NewEncoder(w).Encode(items)
}

func GetOrganization(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var params = mux.Vars(r)
	id, _ := params["id"]

	item := db.GetOrganization(id)

	_ = json.NewEncoder(w).Encode(item)
}

func SendURLTokenForExternalUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var item models.OrganizationForMail
	_ = json.NewDecoder(r.Body).Decode(&item)
	claim := models.External{
		OrganizationID: item.ID,
		Role:           constants.RoleTemp,
	}
	token := mid.GenerateJWTExternal(claim)
	URL := constants.ClientURL + "temp/create-external-user/" + token
	objectMail := models.Mail{
		From: item.Mail,
		Data: URL,
	}
	// aqui debemos validar si el envio del correo se realizo satisfactoriamente o no
	utils.SendMail(objectMail, constants.RouteUserLink)
	_ = json.NewEncoder(w).Encode(URL)
}

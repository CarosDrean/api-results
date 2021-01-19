package controller

import (
	"encoding/json"
	"fmt"
	"github.com/CarosDrean/api-results.git/constants"
	"github.com/CarosDrean/api-results.git/db"
	mid "github.com/CarosDrean/api-results.git/middleware"
	"github.com/CarosDrean/api-results.git/models"
	"github.com/CarosDrean/api-results.git/utils"
	"github.com/gorilla/mux"
	"net/http"
)

type OrganizationController struct {
	DB db.OrganizationDB
}

func (c OrganizationController) GetAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	items, err := c.DB.GetAll()
	if err != nil {
		returnErr(w, err, "obtener todos")
		return
	}
	_ = json.NewEncoder(w).Encode(items)
}

func (c OrganizationController) Get(w http.ResponseWriter, r *http.Request) {
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

func (c OrganizationController) Update(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var item models.Organization
	_ = json.NewDecoder(r.Body).Decode(&item)
	var params = mux.Vars(r)
	id, _ := params["id"]
	result, err := c.DB.Update(id, item)
	if err != nil {
		returnErr(w, err, "update")
		return
	}

	_ = json.NewEncoder(w).Encode(result)
}

func (c OrganizationController) SendURLTokenForExternalUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var item models.OrganizationForMailCreateUser
	_ = json.NewDecoder(r.Body).Decode(&item)
	claim := models.ClaimResult{
		ID:   item.ID,
		Role: constants.Roles.Temp,
		Data: item.TypeUser,
	}

	_ = c.updateURLAdminOrMedic(item)

	token := mid.GenerateJWTExternal(claim)
	URL := constants.ClientURL + "temp/create-external-user/" + token

	organization, _ := c.DB.Get(item.ID)
	objectMail := models.Mail{
		From: item.Mail,
		Data: URL,
		Business: organization.Name,
	}
	data, _ := json.Marshal(objectMail)

	err := utils.SendMail(data, constants.RouteUserLink)
	if err != nil {
		_, _ = fmt.Fprintf(w, "¡Hubo un error al procesar la solicitud!")
		return
	}
	_ = json.NewEncoder(w).Encode(URL)
}

func (c OrganizationController) updateURLAdminOrMedic(item models.OrganizationForMailCreateUser) error  {
	organization, err := c.DB.Get(item.ID)
	if err != nil  {
		return err
	}
	if item.TypeUser == "Admin" {
		organization.UrlAdmin = true
	}
	if item.TypeUser == "Medic" {
		organization.UrlMedic = true
	}
	_, err = c.DB.Update(organization.ID, organization)
	if err != nil {
		return err
	}
	return nil
}

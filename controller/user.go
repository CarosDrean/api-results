package controller

import (
	"encoding/json"
	"github.com/CarosDrean/api-results.git/db"
	"github.com/CarosDrean/api-results.git/models"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func GetSystemUsersPerson(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	res := make([]models.UserPerson, 0)
	items := db.GetSystemUsers()
	for _, e := range items {
		person := db.GetPatient(e.PersonID)[0]
		item := models.UserPerson{
			ID:             e.ID,
			PersonID:       e.PersonID,
			UserName:       e.UserName,
			Password:       e.Password,
			TypeUser:       e.TypeUser,
			OrganizationID: db.GetOrganization(e.OrganizationID).Name,
			DNI:            person.DNI,
			Name:           person.Name,
			FirstLastName:  person.FirstLastName,
			SecondLastName: person.SecondLastName,
			Mail:           person.Mail,
		}
		res = append(res, item)
	}

	_ = json.NewEncoder(w).Encode(res)
}

func GetSystemUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var params = mux.Vars(r)
	id, _ := params["id"]

	items := db.GetSystemUser(id)

	_ = json.NewEncoder(w).Encode(items[0])
}

func UpdatePasswordSystemUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var params = mux.Vars(r)
	id, _ := params["id"]
	var item models.SystemUser
	_ = json.NewDecoder(r.Body).Decode(&item)
	_, err := db.UpdatePasswordSystemUser(id, item.Password)
	if err != nil {
		log.Println(err)
		return
	}

	_ = json.NewEncoder(w).Encode(item)
}

func CreateSystemUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var item models.UserPerson
	_ = json.NewDecoder(r.Body).Decode(&item)
	// verificar si el DNI existe
	person := db.GetPatientFromDNI(item.DNI)
	if len(person) > 0 {
		// solo agregar el system user
	} else {
		// crear la persona y agregar el system user deacuerdo al id
	}
}

func UpdateSystemUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var params = mux.Vars(r)
	id, _ := params["id"]
	var item models.UserPerson
	_ = json.NewDecoder(r.Body).Decode(&item)
	item.ID = id
	person := db.GetPatientFromDNI(item.DNI)
	if len(person) > 0 {
		// solo agregar el system user
	} else {
		// crear la persona y agregar el system user deacuerdo al id
	}
}

func DeleteSystemUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var params = mux.Vars(r)
	id, _ := params["id"]
	result, err := db.DeleteSystemUser(id)
	if err != nil {
		log.Println(err)
	}
	_ = json.NewEncoder(w).Encode(result)
}
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
		person := db.GetPerson(e.PersonID)[0]
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
	personID := validatePerson(item)
	user := models.SystemUser{
		PersonID:       personID,
		UserName:       item.UserName,
		Password:       item.Password,
		TypeUser:       item.TypeUser,
		OrganizationID: item.OrganizationID,
		IsDelete:       0,
	}
	result, err := db.CreateSystemUser(user)
	if err != nil {
		log.Println(err)
	}
	_ = json.NewEncoder(w).Encode(result)
}

// verifica si exste la person y de no ser el caso crea y devuelve el ID
func validatePerson(item models.UserPerson) string {
	personID := item.PersonID
	var err error
	if item.PersonID == "" {
		person := db.GetPersonFromDNI(item.DNI)
		if len(person) == 0 {
			newPerson := models.Person{
				DNI:            item.DNI,
				Password:       "",
				Name:           item.Name,
				FirstLastName:  item.FirstLastName,
				SecondLastName: item.SecondLastName,
				Mail:           item.Mail,
				Sex:            item.Sex,
				Birthday:       item.Birthday,
				IsDeleted:      0,
			}
			personID, err = db.CreatePerson(newPerson)
			if err != nil {
				log.Println(err)
			}
		}
	}
	return personID
}

func UpdateSystemUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var params = mux.Vars(r)
	id, _ := params["id"]
	var item models.UserPerson
	_ = json.NewDecoder(r.Body).Decode(&item)
	item.ID = id

	personId := validatePerson(item)

	user := models.SystemUser{
		PersonID:       personId,
		UserName:       item.UserName,
		Password:       item.Password,
		TypeUser:       item.TypeUser,
		OrganizationID: item.OrganizationID,
		IsDelete:       0,
	}
	result, err := db.UpdateSystemUser(user)
	if err != nil {
		log.Println(err)
	}
	_ = json.NewEncoder(w).Encode(result)
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
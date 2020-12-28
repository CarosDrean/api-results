package controller

import (
	"encoding/json"
	"fmt"
	"github.com/CarosDrean/api-results.git/constants"
	"github.com/CarosDrean/api-results.git/db"
	"github.com/CarosDrean/api-results.git/models"
	"github.com/CarosDrean/api-results.git/utils"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

type UserController struct {}

func (c UserController) GetAllOrganization(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var params = mux.Vars(r)
	idOrganization, _ := params["id"]
	res := make([]models.UserPerson, 0)
	items := db.GetSystemUsers()
	for _, e := range items {
		if e.OrganizationID == idOrganization {
			person := db.GetPerson(e.PersonID)[0]
			item := models.UserPerson{
				ID:             e.ID,
				PersonID:       e.PersonID,
				UserName:       e.UserName,
				Password:       e.Password,
				TypeUser:       e.TypeUser,
				OrganizationID: e.OrganizationID,
				DNI:            person.DNI,
				Name:           person.Name,
				FirstLastName:  person.FirstLastName,
				SecondLastName: person.SecondLastName,
				Mail:           person.Mail,
				Sex:            person.Sex,
				Birthday:       person.Birthday,
			}
			res = append(res, item)
		}
	}

	_ = json.NewEncoder(w).Encode(res)
}

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
			OrganizationID: e.OrganizationID,
			Organization:   db.GetOrganization(e.OrganizationID).Name,
			DNI:            person.DNI,
			Name:           person.Name,
			FirstLastName:  person.FirstLastName,
			SecondLastName: person.SecondLastName,
			Mail:           person.Mail,
			Sex:            person.Sex,
			Birthday:       person.Birthday,
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

	userDB := db.GetSystemUserFromUserName(item.UserName)
	if len(userDB) > 0 {
		_, _ = fmt.Fprintf(w, "¡Nombre de Usuario ya existe en la Base de Datos!")
		return
	}

	personID := validatePerson(item)
	user := models.SystemUser{
		PersonID:       personID,
		UserName:       item.UserName,
		Password:       item.Password,
		TypeUser:       item.TypeUser,
		OrganizationID: item.OrganizationID,
		IsDelete:       0,
	}
	idUser, err := db.CreateSystemUser(user)
	checkError(err, "Created User")

	if user.TypeUser != 1 {
		createProtocolSystemUser(idUser, item.OrganizationID)
	}

	mail := models.Mail{
		From: item.Mail,
		User: item.UserName,
		Password: item.Password,
	}
	utils.SendMail(mail, constants.RouteNewSystemUser)

	_ = json.NewEncoder(w).Encode(idUser)
}

func createProtocolSystemUser(idUser int64, organizationId string) {
	// hay que obtener el id del protocolo deacuerdo a la empresa
	protocols := db.GetProtocolsWidthOrganization(organizationId)
	psu := models.ProtocolSystemUser{
		SystemUserID: idUser,
		ProtocolID:   protocols[0].ID,
	}
	_, err := db.CreateProtocolSystemUser(psu)
	checkError(err, "Created ProtocolSystemUser")
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
			checkError(err, "Created Person")
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
	item.ID, _ = strconv.ParseInt(id, 10, 64)

	personId := validatePerson(item)
	userDB := db.GetSystemUser(strconv.FormatInt(item.ID, 10))
	if userDB[0].UserName != item.UserName {
		userDBo := db.GetSystemUserFromUserName(item.UserName)
		if len(userDBo) > 0 {
			_, _ = fmt.Fprintf(w, "¡Nombre de Usuario ya existe en la Base de Datos!")
			return
		}
	}
	if userDB[0].OrganizationID != item.OrganizationID {
		createProtocolSystemUser(item.ID, item.OrganizationID)
	}

	user := models.SystemUser{
		ID:             item.ID,
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

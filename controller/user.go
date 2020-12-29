package controller

import (
	"encoding/json"
	"fmt"
	"github.com/CarosDrean/api-results.git/constants"
	"github.com/CarosDrean/api-results.git/db"
	"github.com/CarosDrean/api-results.git/models"
	"github.com/CarosDrean/api-results.git/utils"
	"github.com/google/go-cmp/cmp"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

type UserController struct{
	DB db.UserDB
}

func (c UserController) GetAllOrganization(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var params = mux.Vars(r)
	idOrganization, _ := params["id"]
	res := make([]models.UserPerson, 0)
	items, _ := c.DB.GetAll()
	for _, e := range items {
		if e.OrganizationID == idOrganization {
			person, _ := db.PersonDB{}.Get(e.PersonID)
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

func (c UserController) GetAllPerson(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	res := make([]models.UserPerson, 0)
	items, _ := c.DB.GetAll()
	for _, e := range items {
		person, _ := db.PersonDB{}.Get(e.PersonID)
		organization, _ := db.OrganizationDB{}.Get(e.OrganizationID)
		item := models.UserPerson{
			ID:             e.ID,
			PersonID:       e.PersonID,
			UserName:       e.UserName,
			Password:       e.Password,
			TypeUser:       e.TypeUser,
			OrganizationID: e.OrganizationID,
			Organization:   organization.Name,
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

func (c UserController) Get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var params = mux.Vars(r)
	id, _ := params["id"]

	item, _ := c.DB.Get(id)

	_ = json.NewEncoder(w).Encode(item)
}

func (c UserController) UpdatePassword(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var params = mux.Vars(r)
	id, _ := params["id"]
	var item models.SystemUser
	_ = json.NewDecoder(r.Body).Decode(&item)
	_, err := c.DB.UpdatePassword(id, item.Password)
	if err != nil {
		log.Println(err)
		return
	}

	_ = json.NewEncoder(w).Encode(item)
}

func (c UserController) Create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var item models.UserPerson
	_ = json.NewDecoder(r.Body).Decode(&item)

	userDB, _ := c.DB.GetFromUserName(item.UserName)
	if userDB.PersonID != "" && userDB.UserName != "" {
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
	idUser, err := c.DB.Create(user)
	checkError(err, "Created User")

	if user.TypeUser != 1 {
		createProtocolSystemUser(idUser, item.OrganizationID)
	}

	mail := models.Mail{
		From:     item.Mail,
		User:     item.UserName,
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
	newPerson := models.Person{
		DNI:            item.DNI,
		Password:       item.Password,
		Name:           item.Name,
		FirstLastName:  item.FirstLastName,
		SecondLastName: item.SecondLastName,
		Mail:           item.Mail,
		Sex:            item.Sex,
		Birthday:       item.Birthday,
		IsDeleted:      0,
	}
	if item.PersonID == "" {
		person, _ := db.PersonDB{}.GetFromDNI(item.DNI)
		if person.Name == "" && person.DNI == "" {
			personID, err = db.PersonDB{}.Create(newPerson)
			checkError(err, "Created Person")
		}
	} else {
		person, _ := db.PersonDB{}.Get(personID)
		if !cmp.Equal(person, newPerson) {
			fmt.Println("actualizando")
			_, _ = db.PersonDB{}.Update(personID, newPerson)
		}
	}
	return personID
}

func (c UserController) Update(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var params = mux.Vars(r)
	id, _ := params["id"]
	var item models.UserPerson
	_ = json.NewDecoder(r.Body).Decode(&item)
	item.ID, _ = strconv.ParseInt(id, 10, 64)

	personId := validatePerson(item)
	userDB, _ := c.DB.Get(strconv.FormatInt(item.ID, 10))
	if userDB.UserName != item.UserName {
		userDBo, _ := c.DB.GetFromUserName(item.UserName)
		if userDBo.UserName != "" && userDBo.PersonID != ""{
			_, _ = fmt.Fprintf(w, "¡Nombre de Usuario ya existe en la Base de Datos!")
			return
		}
	}
	if userDB.OrganizationID != item.OrganizationID {
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
	result, err := c.DB.Update(user)
	if err != nil {
		log.Println(err)
	}
	_ = json.NewEncoder(w).Encode(result)
}

func (c UserController) Delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var params = mux.Vars(r)
	id, _ := params["id"]
	result, err := c.DB.Delete(id)
	if err != nil {
		log.Println(err)
	}
	_ = json.NewEncoder(w).Encode(result)
}

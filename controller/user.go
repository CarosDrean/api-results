package controller

import (
	"encoding/json"
	"errors"
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
	"time"
)

type UserController struct {
	DB db.UserDB
}

func (c UserController) GetAllOrganization(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var params = mux.Vars(r)
	idOrganization, _ := params["id"]
	res, err := c.usersOrganization(idOrganization)
	if err != nil {
		returnErr(w, err, "obtener todos user")
		return
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
		professional, _ := db.ProfessionalDB{}.Get(e.PersonID)
		item := models.UserPerson{
			ID:               e.ID,
			PersonID:         e.PersonID,
			UserName:         e.UserName,
			Password:         e.Password,
			TypeUser:         e.TypeUser,
			OrganizationID:   e.OrganizationID,
			Organization:     organization.Name,
			DNI:              person.DNI,
			Name:             person.Name,
			FirstLastName:    person.FirstLastName,
			SecondLastName:   person.SecondLastName,
			Mail:             person.Mail,
			Sex:              person.Sex,
			Birthday:         person.Birthday,
			CodeProfessional: professional.Code,
			AccessClient:     e.AccessClient,
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
	professional, _ := db.ProfessionalDB{}.Get(item.PersonID)
	item.CodeProfessional = professional.Code

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

	if !c.permitCreateUser(item.OrganizationID) {
		_, _ = fmt.Fprintf(w, "¡Ya tiene el maximo de usuarios creados!")
		return
	}

	userDB, _ := c.DB.GetFromUserName(item.UserName)
	if userDB.PersonID != "" && userDB.UserName != "" {
		_, _ = fmt.Fprintf(w, "¡Nombre de Usuario ya existe en la Base de Datos!")
		return
	}

	personID, err := c.validateAndCreateOrUpdatePerson(item)
	if err != nil {
		returnErr(w, err, "create")
		return
	}
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
	// 3003 codigo para acceso a clinete
	if user.TypeUser != 1 && user.TypeUser != 5 {
		createProtocolSystemUser(idUser, item.OrganizationID, item.AccessClient)
	}

	mail := models.Mail{
		From:     item.Mail,
		User:     item.UserName,
		Password: item.Password,
	}
	data, _ := json.Marshal(mail)
	_ = utils.SendMail(data, constants.RouteNewSystemUser)

	_ = json.NewEncoder(w).Encode(idUser)
}

func (c UserController) permitCreateUser(idOrganization string) bool {
	if idOrganization == "" {
		return true
	}
	users, _ := c.usersOrganization(idOrganization)
	if len(users) >= constants.MaxUsersOrganization {
		return false
	}
	return true
}

func (c UserController) usersOrganization(idOrganization string) ([]models.UserPerson, error) {
	res := make([]models.UserPerson, 0)
	items, err := c.DB.GetAll()
	if err != nil {
		return res, err
	}
	for _, e := range items {
		if e.OrganizationID == idOrganization && e.TypeUser != constants.CodeRoles.InternalAdmin {
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
	return res, nil
}

func createProtocolSystemUser(idUser int64, organizationId string, accessClient bool) {
	protocols, _ := db.ProtocolDB{}.GetAllOrganization(organizationId)
	psu := models.ProtocolSystemUser{
		SystemUserID: idUser,
		ProtocolID:   protocols[0].ID,
	}
	if accessClient {
		psu.ApplicationHierarchy = constants.CodeAccessClient
	}
	_, err := db.ProtocolSystemUserDB{}.Create(psu)
	checkError(err, "Created ProtocolSystemUser")
}

func updateProtocolSystemUser(idPsu string, idUser int64, organizationId string, accessClient bool) error {
	protocols, _ := db.ProtocolDB{}.GetAllOrganization(organizationId)
	psu := models.ProtocolSystemUser{
		SystemUserID: idUser,
		ProtocolID:   protocols[0].ID,
	}
	if accessClient {
		psu.ApplicationHierarchy = constants.CodeAccessClient
	}
	_, err := db.ProtocolSystemUserDB{}.Update(idPsu, psu)
	checkError(err, "Updated ProtocolSystemUser")
	if err != nil {
		return err
	}
	return nil
}

func (c UserController) validateAndCreateOrUpdateMedic(item models.UserPerson) error {
	personID := item.PersonID
	professional, _ := db.ProfessionalDB{}.Get(personID)
	newProfessional := models.Professional{
		PersonID:     personID,
		ProfessionID: 32, // es medico
		Code:         item.CodeProfessional,
		IsDeleted:    0,
	}
	if professional.Code == "" && professional.ProfessionID == 0 {
		_, err := db.ProfessionalDB{}.Create(newProfessional)
		checkError(err, "Created Professional")
		if err != nil {
			return err
		}
	} else {
		if !cmp.Equal(professional, newProfessional) {
			_, err := db.ProfessionalDB{}.Update(personID, newProfessional)
			checkError(err, "Updated Professional")
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// verifica si exste la person y de no ser el caso crea y devuelve el ID
func (c UserController) validateAndCreateOrUpdatePerson(item models.UserPerson) (string, error) {
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
			item.PersonID = personID
			checkError(err, "Created Person")
		} else {
			item.PersonID = person.ID
		}
	} else {
		person, _ := db.PersonDB{}.Get(personID)
		personCompare := newPerson
		date, _ := time.Parse(time.RFC3339, personCompare.Birthday+"T05:00:00Z")
		personCompare.Birthday = date.String()
		if !cmp.Equal(person, personCompare) {
			fmt.Println("actualizando")
			_, _ = db.PersonDB{}.Update(personID, newPerson)
		}
	}

	if item.TypeUser == constants.CodeRoles.ExternalMedic || item.TypeUser == constants.CodeRoles.ExternalMedicNoData {
		fmt.Println(item.CodeProfessional)
		if item.CodeProfessional == "" {
			return personID, errors.New("code professional invalid")
		}
		err := c.validateAndCreateOrUpdateMedic(item)
		if err != nil {
			return "", err
		}
	}

	return personID, nil
}

func (c UserController) Update(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var params = mux.Vars(r)
	id, _ := params["id"]
	var item models.UserPerson
	_ = json.NewDecoder(r.Body).Decode(&item)
	item.ID, _ = strconv.ParseInt(id, 10, 64)

	personId, err := c.validateAndCreateOrUpdatePerson(item)
	if err != nil {
		returnErr(w, err, "update")
		return
	}
	userDB, _ := c.DB.Get(strconv.FormatInt(item.ID, 10))
	if userDB.UserName != item.UserName {
		userDBo, _ := c.DB.GetFromUserName(item.UserName)
		if userDBo.UserName != "" && userDBo.PersonID != "" {
			_, _ = fmt.Fprintf(w, "¡Nombre de Usuario ya existe en la Base de Datos!")
			return
		}
	}
	psuDB, _ := db.ProtocolSystemUserDB{}.GetAllSystemUserID(strconv.Itoa(int(userDB.ID)))
	if userDB.OrganizationID != item.OrganizationID || !(psuDB[0].ApplicationHierarchy == constants.CodeAccessClient && item.AccessClient) {
		err = updateProtocolSystemUser(psuDB[0].ID, item.ID, item.OrganizationID, item.AccessClient)
		if err != nil {
			returnErr(w, err, "update protocol system user")
		}
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
		returnErr(w, err, "update")
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

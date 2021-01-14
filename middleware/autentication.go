package middleware

import (
	"encoding/json"
	"fmt"
	"github.com/CarosDrean/api-results.git/constants"
	"github.com/CarosDrean/api-results.git/db"
	"github.com/CarosDrean/api-results.git/helper"
	"github.com/CarosDrean/api-results.git/models"
	"log"
	"net/http"
)

func Login(w http.ResponseWriter, r *http.Request){
	var user models.UserLogin
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		_, _ = fmt.Fprintf(w, "Error al leer el usuario %s\n", err)
		return
	}
	log.Println(user)
	var stateLogin constants.State
	var id string
	isSystemUser := false
	isPatientParticular := false
	if !user.Particular {
		stateLogin, id, err = patientBusiness(user)

		if stateLogin == constants.NotFound {
			stateLogin, id, err = db.UserDB{}.ValidateLogin(user.User, user.Password)
			isSystemUser = true
		}
	} else {
		isPatientParticular = true
		isSystemUser = false
		stateLogin, id, err = patientParticular(user)
	}
	if err != nil {
		_, _ = fmt.Fprintf(w, fmt.Sprintf("¡Hubo un Error %s", err.Error()))
	}


	switch stateLogin {
	case constants.Accept:
		userResult := models.ClaimResult{ID: id, Role: getRole(0)}
		if isSystemUser {
			systemUser, _ := db.UserDB{}.Get(id)
			userResult = models.ClaimResult{ID: id, Role: getRole(systemUser.TypeUser)}
		}
		token := GenerateJWT(userResult)
		result := models.ResponseToken{Token: token}
		jsonResult, err := json.Marshal(result)
		if err != nil {
			fmt.Println(w, "Error al generar el json")
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(jsonResult)
	case constants.ErrorUP:
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprintf(w, "Hubo un error!")
		break
	case constants.NotFoundMail:
		w.WriteHeader(http.StatusForbidden)
		_, _ = fmt.Fprintf(w, "No se encontro su direccion de correo electronico")
		break
	case constants.NotFound:
		w.WriteHeader(http.StatusForbidden)
		if !isSystemUser || isPatientParticular{
			_, _ = fmt.Fprintf(w, "¡No existe Paciente!")
		} else {
			_, _ = fmt.Fprintf(w, "¡No existe Usuario!")
		}
	case constants.InvalidCredentials:
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = fmt.Fprintf(w, "¡Contraseña Incorrecta!")
		break
	case constants.PasswordUpdate:
		w.WriteHeader(http.StatusFound)
		_, _ = fmt.Fprintf(w, "Consulte su correo electronico con las nuevas credenciales.")
		break
	}
}

// en las dos funciones siguientes inicializamos la bd dependiendo del caso, tambien lo dejamos en el main por si acaso
func patientBusiness(user models.UserLogin) (constants.State, string, error){
	db.DB = helper.Get()
	return db.PersonDB{}.ValidateLogin(user.User, user.Password)
}

func patientParticular(user models.UserLogin) (constants.State, string, error){
	db.DB = helper.GetAux()
	return db.PersonDB{}.ValidateLogin(user.User, user.Password)
}

func getRole(typeUser int)constants.Role{
	switch typeUser {
	case constants.CodeRoles.Patient:
		return constants.Roles.Patient
	case constants.CodeRoles.InternalAdmin:
		return constants.Roles.InternalAdmin
	case constants.CodeRoles.ExternalAdmin:
		return constants.Roles.ExternalAdmin
	case constants.CodeRoles.ExternalMedic:
		return constants.Roles.ExternalMedic
	case constants.CodeRoles.ExternalMedicNoData:
		return constants.Roles.ExternalMedic
	default:
		return ""
	}
}

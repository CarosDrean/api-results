package controller

import (
	_ "context"
	_ "database/sql"
	"encoding/json"
	_ "encoding/json"
	_ "errors"
	_ "fmt"
	_ "github.com/CarosDrean/api-results.git/constants"
	"github.com/CarosDrean/api-results.git/db"
	_ "github.com/CarosDrean/api-results.git/db"
	"github.com/CarosDrean/api-results.git/models"
	_ "github.com/CarosDrean/api-results.git/models"
	_ "github.com/CarosDrean/api-results.git/query"
	_ "github.com/CarosDrean/api-results.git/utils"
	_ "github.com/google/go-cmp/cmp"
	_ "github.com/gorilla/mux"
	_ "log"
	"net/http"
	_ "net/http"
	_ "strconv"
	_ "time"
)

type PetitionController struct {
	DB db.PetitionDB
	DBB db.CitaDB
}

func (c PetitionController) Create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var item models.PetitionProgrammation
	_ = json.NewDecoder(r.Body).Decode(&item)

	petit := models.PetitionProgrammation{
		PersonId:       		item.PersonId,
		DocType:        		item.DocType,
		DocNumber:       		item.DocNumber,
		FirstName:       		item.FirstName,
		FirstLastName: 			item.FirstLastName,
		SecondLastName:     	item.SecondLastName,
		SexTypeId: 				item.SexTypeId,
		Birthdate: 				item.Birthdate,
		TelephoneNumber: 		item.TelephoneNumber,
		CurrentOccupation: 		item.CurrentOccupation,
		DateProgramming: 		item.DateProgramming,
		ServiceTypeId: 			item.ServiceTypeId,
		PersonProgramming: 		item.PersonProgramming,
		ResponsableProgramming: item.ResponsableProgramming,
		CalendarId_2: 			item.CalendarId_2,
		WorkersCondition: 		item.WorkersCondition,
		FactCR:				 	item.FactCR,
		NombreProyecto: 		item.NombreProyecto,
		OrganizationId: 		item.OrganizationId,
		ProtocolId: 			item.ProtocolId,
		Deleted: 				item.Deleted,
		PetitionStatus: 		item.PetitionStatus,
		Comentary: 				item.Comentary,
	}

	result, err := c.DB.Create(petit)
	if err != nil {
		returnErr(w, err, "createdPetition")
		return
	}

	_ = json.NewEncoder(w).Encode(result)
}

func (c PetitionController) CreateCita(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var item models.MailConsultaCardiologica
	_ = json.NewDecoder(r.Body).Decode(&item)

	cita := models.MailConsultaCardiologica{
		Nombre: item.Nombre,
		Apepaterno: item.Apepaterno,
		Apematerno: item.Apematerno,
		Dni: item.Dni,
		Email: item.Email,
		Telefono: item.Telefono,
		Direccion: item.Direccion,
		Dob: item.Dob,
		Fecha: item.Fecha,
		Sexo: item.Sexo,
		Mensaje: item.Mensaje,
	}

	result, err := c.DBB.CreateCita(cita)
	if err != nil {
		returnErr(w, err, "createdCita")
		return
	}

	_ = json.NewEncoder(w).Encode(result)
}
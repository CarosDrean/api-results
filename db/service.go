package db

import (
	"database/sql"
	"fmt"
	"github.com/CarosDrean/api-results.git/constants"
	"github.com/CarosDrean/api-results.git/models"
	"github.com/CarosDrean/api-results.git/query"
	"strings"
)

type ServiceDB struct{}

func (db ServiceDB) GetAllCovid(docNumber string) ([]models.ServiceCovid, error) {
	res := make([]models.ServiceCovid, 0)
	var item models.ServiceCovid
	tsql := fmt.Sprintf(query.Service["getAllCovid"].Q, docNumber)
	rows, err := DB.Query(tsql)

	if err != nil {
		fmt.Println("Error reading rows: " + err.Error())
		return res, err
	}
	for rows.Next() {
		err := rows.Scan(&item.Date, &item.Name, &item.FirstLastname, &item.SecondLastName, &item.DocNumber, &item.BirthDate,
			&item.Sex, &item.Group, &item.Occupation, &item.Exam, &item.Result)
		if err != nil {
			checkError(err, "next", "db", "getallcovid")
		} else {
			// calcular edad
			item.Age = 30
			res = append(res, item)
		}
	}
	return res, nil
}

func (db ServiceDB) GetAllDate(filter models.Filter) ([]models.ServicePatientOrganization, error) {
	res := make([]models.ServicePatientOrganization, 0)
	var service models.Service
	var person models.Person
	var organization models.Organization
	var protocol models.Protocol

	tsql := fmt.Sprintf(query.Service["listDate"].Q, filter.DateFrom, filter.DateTo)
	rows, err := DB.Query(tsql)

	if err != nil {
		fmt.Println("Error reading rows: " + err.Error())
		return res, err
	}
	for rows.Next() {
		var pass sql.NullString
		var birth sql.NullString
		var phone sql.NullString
		var occupation sql.NullString
		var doc sql.NullInt64
		err := rows.Scan(&service.ID, &service.PersonID, &service.ProtocolID, &service.ServiceDate, &service.ServiceStatusId,
			&service.IsDeleted, &service.AptitudeStatusId,
			&person.ID, &person.DNI, &pass, &person.Name, &person.FirstLastName, &person.SecondLastName, &person.Mail,
			&person.Sex, &birth, &person.IsDeleted, &phone, &occupation, &doc,
			&organization.ID, &organization.Name,
			&protocol.EsoType)
		person.Password = pass.String
		person.Birthday = birth.String
		person.Phone = phone.String
		if err != nil {
			checkError(err, "next", "Db", "GetAlDate")
		} else if service.IsDeleted != 1 && service.ServiceStatusId == 3 {

			result, _ := ResultDB{}.GetService(service.ID, constants.IdPruebaRapida, constants.IdResultPruebaRapida)
			result2, _ := ResultDB{}.GetService(service.ID, constants.IdPruebaHisopado, constants.IdResultPruebaHisopado)
			iStatusliquidation, _ := StatusGenerateBD{}.GetStatusGenerate(service.ID)
			item := models.ServicePatientOrganization{
				ID:               service.ID,
				ServiceDate:      service.ServiceDate,
				PersonID:         service.PersonID,
				ProtocolID:       service.ProtocolID,
				OrganizationID:   organization.ID,
				Organization:     organization.Name,
				AptitudeStatusId: service.AptitudeStatusId,
				DNI:              person.DNI,
				Name:             person.Name,
				FirstLastName:    person.FirstLastName,
				SecondLastName:   person.SecondLastName,
				Mail:             person.Mail,
				Sex:              person.Sex,
				Birthday:         person.Birthday,
				EsoType:          protocol.EsoType,
				Phone:            person.Phone,
				Occupation:       person.Occupation,
				Doc:              person.Doc,
				Result:           result,
				Result2:          result2,
				GenerateStatus:   iStatusliquidation,
			}
			res = append(res, item)
		}
	}
	defer rows.Close()
	return res, nil
}

func (db ServiceDB) GetAllExamDetail(filter models.Filter) ([]models.ServicePatientExam, error) {
	res := make([]models.ServicePatientExam, 0)
	// var item models.ServicePatientExam
	tsql := fmt.Sprintf(query.Service["listExamsDetailDate"].Q, filter.DateFrom, filter.DateTo)
	rows, err := DB.Query(tsql)

	if err != nil {
		fmt.Println("Error reading rows: " + err.Error())
		return res, err
	}
	/*for rows.Next() {
		var birth sql.NullString
		err := rows.Scan(&item.ID, &item.ProtocolID, &item.LocationID, &item.OrganizationID, &item.FirstLastName, &item.SecondLastName, &item.Name,
			&item.TypeDoc, &item.DNI, &item.Organization, &item.Occupation, &item.ServiceDate, &birth, &item.PriceProtocol,
			&item.Component, &item.PriceComponent, &item.Protocol, &item.Mail, &item.Sex, &item.AptitudeStatusId, &item.EsoType)
		item.Birthday = birth.String
		if err != nil {
			log.Println(err)
		} else {
			res = append(res, item)
		}
	}*/
	defer rows.Close()
	return res, nil
}

func (db ServiceDB) GetAllPerson(id string) ([]models.Service, error) {
	res := make([]models.Service, 0)

	tsql := fmt.Sprintf(query.Service["getPersonID"].Q, id)
	rows, err := DB.Query(tsql)

	err = db.scan(rows, err, &res, "Service DB", "GetAll Person")
	if err != nil {
		return res, err
	}
	defer rows.Close()
	return res, nil
}

func (db ServiceDB) GetAllProtocol(id string) ([]models.Service, error) {
	res := make([]models.Service, 0)

	tsql := fmt.Sprintf(query.Service["getProtocol"].Q, id)
	rows, err := DB.Query(tsql)

	err = db.scan(rows, err, &res, "Service DB", "GetAll Person")
	if err != nil {
		return res, err
	}
	defer rows.Close()
	return res, nil
}

func (db ServiceDB) GetAllDiseaseFilterDate(filter models.Filter) []models.ServicePatientDiseases {
	res := make([]models.ServicePatientDiseases, 0)
	var service models.Service
	var person models.Person
	var protocol models.Protocol

	tsql := fmt.Sprintf(query.Service["listDiseaseFilter"].Q, filter.DateFrom, filter.DateTo)
	rows, err := DB.Query(tsql)

	if err != nil {
		fmt.Println("Error reading rows: " + err.Error())
		return res
	}
	for rows.Next() {
		var pass sql.NullString
		var birth sql.NullString
		var disease sql.NullString
		var phone sql.NullString
		var occupation sql.NullString
		var doc sql.NullInt64
		var mail sql.NullString
		var diseaseString string
		err := rows.Scan(&service.ID, &service.PersonID, &service.ProtocolID, &service.ServiceDate, &service.ServiceStatusId,
			&service.IsDeleted, &service.AptitudeStatusId,
			&person.ID, &person.DNI, &pass, &person.Name, &person.FirstLastName, &person.SecondLastName, &mail,
			&person.Sex, &birth, &person.IsDeleted, &phone, &occupation, &doc,
			&protocol.ID, &protocol.Name, &protocol.OrganizationID, &protocol.OrganizationEmployerID, &protocol.LocationID, &protocol.EsoType,
			&protocol.GroupOccupationId,
			&disease)
		if pass.Valid {
			person.Password = pass.String
		} else {
			person.Password = ""
		}
		if birth.Valid {
			person.Birthday = birth.String
		} else {
			person.Birthday = ""
		}
		if disease.Valid {
			diseaseString = disease.String
		} else {
			diseaseString = ""
		}
		if err != nil {
			checkError(err, "next", "db", "get all disease filter")
		} else if service.IsDeleted != 1 {
			item := models.ServicePatientDiseases{
				ID:               service.ID,
				ServiceDate:      service.ServiceDate,
				PersonID:         service.PersonID,
				ProtocolID:       service.ProtocolID,
				OrganizationID:   protocol.OrganizationID,
				AptitudeStatusId: service.AptitudeStatusId,
				DNI:              person.DNI,
				Name:             person.Name,
				FirstLastName:    person.FirstLastName,
				SecondLastName:   person.SecondLastName,
				Mail:             person.Mail,
				Sex:              person.Sex,
				Birthday:         person.Birthday,
				Phone:            person.Phone,
				Occupation:       person.Occupation,
				Doc:              person.Doc,
				Disease:          diseaseString,
				EsoType:          protocol.EsoType,
			}
			res = append(res, item)
		}
	}
	defer rows.Close()
	return res
}

func (db ServiceDB) GetAll() ([]models.Service, error) {
	res := make([]models.Service, 0)

	tsql := fmt.Sprintf(query.Service["list"].Q)
	rows, err := DB.Query(tsql)

	err = db.scan(rows, err, &res, "Service DB", "GetAll")
	if err != nil {
		return res, err
	}
	defer rows.Close()
	return res, nil
}

func (db ServiceDB) Get(id string) (models.Service, error) {
	res := make([]models.Service, 0)

	tsql := fmt.Sprintf(query.Service["get"].Q, id)
	rows, err := DB.Query(tsql)

	err = db.scan(rows, err, &res, "Service DB", "Get")
	if err != nil {
		return models.Service{}, err
	}
	if len(res) == 0 {
		return models.Service{}, nil
	}
	defer rows.Close()
	return res[0], nil
}

func (db ServiceDB) GetAllProtocolFilter(id string, filter models.Filter) ([]models.Service, error) {
	res := make([]models.Service, 0)

	tsql := fmt.Sprintf(query.Service["getProtocolFilter"].Q, id, filter.DateFrom, filter.DateTo)
	rows, err := DB.Query(tsql)

	err = db.scan(rows, err, &res, "Service DB", "GetAll protocol")
	if err != nil {
		return res, err
	}
	defer rows.Close()
	return res, nil
}

// deacuerdo al id de la empresa obtiene todos sus protocolos y va armando el objeto ServicePatient
// obtener las empresas
func (db ServiceDB) GetAllPatientsWithOrganizationFilter(idOrganization string, filter models.Filter) ([]models.ServicePatient, error) {
	res := make([]models.ServicePatient, 0)

	// esto esta mal, deberia obtener los pacientes defrente o los servicios
	protocols, err := ProtocolDB{}.GetAllOrganization(idOrganization)
	if err != nil {
		return res, err
	}
	for _, e := range protocols {
		sp, _ := db.GetAllPatientsWithProtocolFilter(e.ID, filter, false)
		res = append(res, sp...)
	}
	return res, nil
}

// obtener las empresas con su contratista
func (db ServiceDB) GetAllPatientsWithOrganizationEmployerFilter(filter models.Filter) ([]models.ServicePatient, error) {
	res := make([]models.ServicePatient, 0)

	protocols, err := ProtocolDB{}.GetAllOrganizationEmployer(filter.ID)
	if err != nil {
		return res, err
	}
	for _, e := range protocols {
		sp, _ := db.GetAllPatientsWithProtocolFilter(e.ID, filter, true)
		res = append(res, sp...)
	}
	return res, nil
}

func (db ServiceDB) GetAllPatientsWithLocationFilter(idLocation string, filter models.Filter) ([]models.ServicePatient, error) {
	res := make([]models.ServicePatient, 0)

	protocols, err := ProtocolDB{}.GetAllLocation(idLocation)
	if err != nil {
		return res, err
	}
	for _, e := range protocols {
		sp, _ := db.GetAllPatientsWithProtocolFilter(e.ID, filter, false)
		res = append(res, sp...)
	}
	return res, nil
}

func (db ServiceDB) GetAllPatientsWithProtocolFilter(idProtocol string, filter models.Filter, needOrganization bool) ([]models.ServicePatient, error) {
	res := make([]models.ServicePatient, 0)

	services, err := db.GetAllProtocolFilter(idProtocol, filter)
	if err != nil {
		return res, err
	}
	for _, e := range services {
		patient, _ := PersonDB{}.Get(e.PersonID)
		result, _ := ResultDB{}.GetService(e.ID, constants.IdPruebaRapida, constants.IdResultPruebaRapida)
		result2, _ := ResultDB{}.GetService(e.ID, constants.IdPruebaHisopado, constants.IdResultPruebaHisopado)
		iStatusliquidation, _ := StatusGenerateBD{}.GetStatusGenerate(e.ID)
		calendar, _ := CalendarDB{}.GetService(e.ID)
		// para quitar la zona horaria
		start := strings.Split(calendar.CircuitStart, ".")
		end := strings.Split(calendar.CircuitEnd, ".")
		item := models.ServicePatient{
			ID:               e.ID,
			ServiceDate:      e.ServiceDate,
			PersonID:         patient.ID,
			ProtocolID:       e.ProtocolID,
			Birthday:         patient.Birthday,
			AptitudeStatusId: e.AptitudeStatusId,
			DNI:              patient.DNI,
			Name:             patient.Name,
			FirstLastName:    patient.FirstLastName,
			SecondLastName:   patient.SecondLastName,
			Mail:             patient.Mail,
			Sex:              patient.Sex,
			Result:           result,
			Result2:          result2,
			CalendarStatus:   calendar.CalendarStatusID,
			CircuitStart:     start[0],
			CircuitEnd:       end[0],
			GenerateStatus:   iStatusliquidation,
		}
		if needOrganization {
			protocol, _ := ProtocolDB{}.Get(e.ProtocolID)
			organization, _ := OrganizationDB{}.Get(protocol.OrganizationID)

			item.OrganizationID = organization.ID
			item.Organization = organization.Name

		}
		res = append(res, item)
	}
	return res, nil
}

func (db ServiceDB) GetGesoAuxFilter(idOrganization string, filter models.Filter) ([]models.ServicePatient, error) {
	res := make([]models.ServicePatient, 0)
	var item models.ServicePatient
	tsql := fmt.Sprintf(query.Service["getGesoFilter"].Q, idOrganization, filter.DateFrom, filter.DateTo)
	rows, err := DB.Query(tsql)

	if err != nil {
		checkError(err, "GetGesoFilter", "db", "Reading rows")
		return res, err
	}

	for rows.Next() {
		err = rows.Scan(&item.ID, &item.ServiceDate, &item.PersonID, &item.ProtocolID, &item.AptitudeStatusId, &item.DNI,
			&item.Name, &item.FirstLastName, &item.SecondLastName, &item.Mail, &item.Sex, &item.Birthday, &item.CalendarStatus, &item.CircuitStart,
			&item.CircuitEnd, &item.Geso)
		if err != nil {
			checkError(err, "GetGesoFilter", "db", "scan rows")
		} else {
			res = append(res, item)
		}

	}
	return res, nil
}

func (db ServiceDB) GetGesoFilter(idOrganization string, filter models.Filter) ([]models.ServicePatient, error) {
	res := make([]models.ServicePatient, 0)
	//var items models.ServicePatient

	services, err := db.GetGesoAuxFilter(idOrganization, filter)
	if err != nil {
		return res, err
	}
	for _, e := range services {
		patient, _ := PersonDB{}.Get(e.PersonID)
		result, _ := ResultDB{}.GetService(e.ID, constants.IdPruebaRapida, constants.IdResultPruebaRapida)
		result2, _ := ResultDB{}.GetService(e.ID, constants.IdPruebaHisopado, constants.IdResultPruebaHisopado)
		iStatusliquidation, _ := StatusGenerateBD{}.GetStatusGenerate(e.ID)
		calendar, _ := CalendarDB{}.GetService(e.ID)
		// para quitar la zona horaria
		start := strings.Split(calendar.CircuitStart, ".")
		end := strings.Split(calendar.CircuitEnd, ".")
		item := models.ServicePatient{
			ID:               e.ID,
			ServiceDate:      e.ServiceDate,
			PersonID:         patient.ID,
			ProtocolID:       e.ProtocolID,
			Birthday:         patient.Birthday,
			AptitudeStatusId: e.AptitudeStatusId,
			DNI:              patient.DNI,
			Name:             patient.Name,
			FirstLastName:    patient.FirstLastName,
			SecondLastName:   patient.SecondLastName,
			Mail:             patient.Mail,
			Sex:              patient.Sex,
			Result:           result,
			Result2:          result2,
			GenerateStatus:   iStatusliquidation,
			CalendarStatus:   calendar.CalendarStatusID,
			CircuitStart:     start[0],
			CircuitEnd:       end[0],
			Geso:             e.Geso,
		}

		res = append(res, item)
	}
	return res, nil
}

func (db ServiceDB) scan(rows *sql.Rows, err error, res *[]models.Service, ctx string, situation string) error {
	var item models.Service
	if err != nil {
		checkError(err, situation, ctx, "Reading rows")
		return err
	}
	for rows.Next() {
		err := rows.Scan(&item.ID, &item.PersonID, &item.ProtocolID, &item.ServiceDate, &item.ServiceStatusId,
			&item.IsDeleted, &item.AptitudeStatusId)
		if err != nil {
			checkError(err, situation, ctx, "Scan rows")
			return err
		} else if item.IsDeleted != 1 /*&& item.ServiceStatusId == 3*/ { // verificar servicios no eliminados y culminados
			*res = append(*res, item)
		}
	}
	return nil
}

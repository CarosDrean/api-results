package db

import (
	"database/sql"
	"fmt"
	"github.com/CarosDrean/api-results.git/models"
	"github.com/CarosDrean/api-results.git/query"
	"log"
)

type ServiceDB struct {}

func (db ServiceDB) GetAllPerson(id string) ([]models.Service, error)  {
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

func (db ServiceDB) GetAllProtocol(id string) ([]models.Service, error)  {
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
		var diseaseString string
		err := rows.Scan(&service.ID, &service.PersonID, &service.ProtocolID, &service.ServiceDate, &service.ServiceStatusId,
			&service.IsDeleted, &service.AptitudeStatusId,
			&person.ID, &person.DNI, &pass, &person.Name, &person.FirstLastName, &person.SecondLastName, &person.Mail,
			&person.Sex, &birth, &person.IsDeleted,
			&protocol.ID, &protocol.Name, &protocol.OrganizationID, &protocol.LocationID, &protocol.IsDeleted, &protocol.EsoType,
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
			log.Println(err)
		} else {
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
		} else if item.IsDeleted != 1 && item.ServiceStatusId == 3 { // verificar servicios no eliminados y culminados
			*res = append(*res, item)
		}
	}
	return nil
}


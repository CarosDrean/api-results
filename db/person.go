package db

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/CarosDrean/api-results.git/constants"
	"github.com/CarosDrean/api-results.git/models"
	"github.com/CarosDrean/api-results.git/query"
	"github.com/CarosDrean/api-results.git/utils"
	"strings"
	"time"
)

type PersonDB struct {}

func (db PersonDB) Get(id string) (models.Person, error) {
	res := make([]models.Person, 0)

	tsql := fmt.Sprintf(query.Person["get"].Q, id)
	rows, err := DB.Query(tsql)

	err = db.scan(rows, err, &res, "Person DB", "Get")
	if err != nil {
		return models.Person{}, err
	}
	if len(res) == 0 {
		return models.Person{}, nil
	}
	defer rows.Close()
	return res[0], nil
}

func (db PersonDB) GetFromDNI(dni string) (models.Person, error) {
	res := make([]models.Person, 0)

	tsql := fmt.Sprintf(query.Person["getDNI"].Q, dni)
	rows, err := DB.Query(tsql)

	err = db.scan(rows, err, &res, "Person DB", "Get")
	if err != nil {
		return models.Person{}, err
	}
	if len(res) == 0 {
		return models.Person{}, nil
	}
	defer rows.Close()
	return res[0], nil
}

func (db PersonDB) Create(item models.Person) (string, error) {
	ctx := context.Background()
	tsql := fmt.Sprintf(query.Person["insert"].Q)
	if item.Password != "" {
		item.Password = encryptMD5(item.Password)
	}
	sqdb :=SequentialDB{}
	sequentialID := sqdb.NextSequentialId(constants.IdNode, constants.IdPersonTable)
	newId := sqdb.NewID(constants.IdNode, sequentialID, constants.PrefixPerson)
	item.ID = newId

	date, _ := time.Parse(time.RFC3339, item.Birthday + "T05:00:00Z")

	_, err := DB.ExecContext(
		ctx,
		tsql,
		sql.Named("v_PersonId", item.ID),
		sql.Named("v_DocNumber", item.DNI),
		sql.Named("v_Password", item.Password),
		sql.Named("v_FirstName", item.Name),
		sql.Named("v_FirstLastName", item.FirstLastName),
		sql.Named("v_SecondLastName", item.SecondLastName),
		sql.Named("v_Mail", item.Mail),
		sql.Named("i_SexTypeId", item.Sex),
		sql.Named("d_Birthdate", date),
		sql.Named("i_IsDeleted", 0))
	if err != nil {
		return "", err
	}
	return newId, nil
}

func (db PersonDB) Update(id string, item models.Person) (int64, error) {
	ctx := context.Background()
	tsql := fmt.Sprintf(query.Person["update"].Q)

	date, _ := time.Parse(time.RFC3339, item.Birthday + "T05:00:00Z")
	result, err := DB.ExecContext(
		ctx,
		tsql,
		sql.Named("ID", id),
		sql.Named("v_PersonId", item.ID),
		sql.Named("v_DocNumber", item.DNI),
		sql.Named("v_Password", item.Password),
		sql.Named("v_FirstName", item.Name),
		sql.Named("v_FirstLastName", item.FirstLastName),
		sql.Named("v_SecondLastName", item.SecondLastName),
		sql.Named("v_Mail", item.Mail),
		sql.Named("i_SexTypeId", item.Sex),
		sql.Named("d_Birthdate", date),
		sql.Named("i_IsDeleted", 0))
	if err != nil {
		return -1, err
	}
	return result.RowsAffected()
}

func (db PersonDB) ValidateLogin(user string, password string) (constants.State, string){
	item, err := db.GetFromDNI(user)
	if err != nil {
		return constants.NotFound, ""
	}
	if item.DNI == "" && item.Name == "" {
		return constants.NotFound, ""
	}
	if validatePasswordPatientForReset(password, item){
		if len(item.Mail) != 0{
			newPassword := utils.CreateNewPassword()
			mail := models.Mail{
				From: item.Mail,
				User: user,
				Password: newPassword,
			}
			_, err := db.UpdatePassword(item.ID, newPassword)
			if err != nil {
				return constants.ErrorUP, ""
			}
			_ = utils.SendMail(mail, constants.RouteNewPassword)
			return constants.PasswordUpdate, ""
		}
		return constants.NotFoundMail, ""
	}
	if item.Password == password {
		return constants.Accept, item.ID
	}
	return constants.InvalidCredentials, ""
}

func (db PersonDB) UpdatePassword(id string, password string) (int64, error) {
	ctx := context.Background()
	tsql := fmt.Sprintf(query.Person["updatePassword"].Q, id)
	result, err := DB.ExecContext(
		ctx,
		tsql,
		sql.Named("Password", password))
	if err != nil {
		return -1, err
	}
	return result.RowsAffected()
}

func validatePasswordPatientForReset(password string, patient models.Person) bool {
	return patient.DNI == password
}

func (db PersonDB) scan(rows *sql.Rows, err error, res *[]models.Person, ctx string, situation string) error {
	var item models.Person
	if err != nil {
		checkError(err, situation, ctx, "Reading rows")
		return err
	}
	for rows.Next() {
		var pass sql.NullString
		var birth sql.NullString
		err := rows.Scan(&item.ID, &item.DNI, &pass, &item.Name, &item.FirstLastName, &item.SecondLastName, &item.Mail,
			&item.Sex, &birth, &item.IsDeleted)
		item.Birthday = birth.String
		if pass.Valid {
			item.Password = pass.String
		} else {
			item.Password = item.DNI
		}
		if strings.Contains(item.Mail, "notiene") || strings.Contains(item.Mail, "NOTIENE"){
			item.Mail = ""
		}
		if err != nil {
			checkError(err, situation, ctx, "Scan rows")
			return err
		} else {
			*res = append(*res, item)
		}
	}
	return nil
}


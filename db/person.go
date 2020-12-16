package db

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/CarosDrean/api-results.git/constants"
	"github.com/CarosDrean/api-results.git/models"
	"github.com/CarosDrean/api-results.git/utils"
	"log"
	"strings"
)

func GetPerson(id string) []models.Person {
	res := make([]models.Person, 0)
	var item models.Person

	tsql := fmt.Sprintf(QueryPerson["get"].Q, id)
	rows, err := DB.Query(tsql)

	if err != nil {
		fmt.Println("Error reading rows: " + err.Error())
		return res
	}
	for rows.Next(){
		var pass sql.NullString
		var birth sql.NullString
		err := rows.Scan(&item.ID, &item.DNI, &pass, &item.Name, &item.FirstLastName, &item.SecondLastName, &item.Mail,
			&item.Sex, &birth, &item.IsDeleted)
		if pass.Valid {
			item.Password = pass.String
		} else {
			item.Password = ""
		}
		if birth.Valid {
			item.Birthday = birth.String
		} else {
			item.Birthday = ""
		}
		if err != nil {
			log.Println(err)
			return res
		} else{
			res = append(res, item)
		}
 	}
	defer rows.Close()
	return res
}

func GetPersonFromDNI(dni string) []models.Person {
	res := make([]models.Person, 0)
	var item models.Person

	tsql := fmt.Sprintf(QueryPerson["getDNI"].Q, dni)
	rows, err := DB.Query(tsql)

	if err != nil {
		fmt.Println("Error reading rows: " + err.Error())
		return res
	}
	for rows.Next(){
		var pass sql.NullString
		err := rows.Scan(&item.ID, &item.DNI, &pass, &item.Name, &item.FirstLastName, &item.SecondLastName, &item.Mail,
			&item.Sex, &item.Birthday, &item.IsDeleted)
		if pass.Valid {
			item.Password = pass.String
		} else {
			item.Password = dni
		}
		if strings.Contains(item.Mail, "notiene") || strings.Contains(item.Mail, "NOTIENE"){
			item.Mail = ""
		}
		if err != nil {
			log.Println("error...........")
			log.Println(err)
			return res
		} else{
			res = append(res, item)
			log.Println(item.Password)
		}
	}
	defer rows.Close()
	return res
}

func CreatePerson(item models.Person) (string, error) {
	ctx := context.Background()
	tsql := fmt.Sprintf(QueryPerson["insert"].Q)
	if item.Password != "" {
		item.Password = encryptMD5(item.Password)
	}

	sequentialID := GetNextSequentialId(constants.IdNode, constants.IdPersonTable)
	newId := GetNewID(constants.IdNode, sequentialID, constants.PrefixPerson)
	item.ID = newId

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
		sql.Named("d_Birthdate", item.Birthday),
		sql.Named("i_IsDeleted", 0))
	if err != nil {
		return "", err
	}
	return newId, nil
}

func ValidatePatientLogin(user string, password string) (constants.State, string){
	items := GetPersonFromDNI(user)
	if len(items) > 0 {
		if validatePasswordPatientForReset(password, items[0]){
			if len(items[0].Mail) != 0{
				newPassword := utils.CreateNewPassword()
				mail := models.Mail{
					From: items[0].Mail,
					User: user,
					Password: newPassword,
				}
				_, err := UpdatePasswordPatient(items[0].ID, newPassword)
				if err != nil {
					return constants.ErrorUP, ""
				}
				utils.Sendmail(mail)
				return constants.PasswordUpdate, ""
			}
			return constants.NotFoundMail, ""
		}
		if items[0].Password == password {
			return constants.Accept, items[0].ID
		}
		return constants.InvalidCredentials, ""
	}
	return constants.NotFound, ""
}

func UpdatePasswordPatient(id string, password string) (int64, error) {
	ctx := context.Background()
	tsql := fmt.Sprintf(QueryPerson["updatePassword"].Q, id)
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


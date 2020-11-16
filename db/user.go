package db

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/CarosDrean/api-results.git/helper"
	"github.com/CarosDrean/api-results.git/models"
	"github.com/CarosDrean/api-results.git/utils"
	"golang.org/x/crypto/bcrypt"
	"log"
)

func GetSystemUser(id string) []models.SystemUser {
	res := make([]models.SystemUser, 0)
	var item models.SystemUser

	tsql := fmt.Sprintf(QuerySystemUser["get"].Q, id)
	rows, err := DB.Query(tsql)

	if err != nil {
		fmt.Println("Error reading rows: " + err.Error())
		return res
	}
	for rows.Next(){
		err := rows.Scan(&item.ID, &item.PersonID, &item.UserName, &item.Password, &item.TypeUser)
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

func GetSystemUserFromUserName(userName string) []models.SystemUser {
	res := make([]models.SystemUser, 0)
	var item models.SystemUser

	tsql := fmt.Sprintf(QuerySystemUser["getUserName"].Q, userName)
	rows, err := DB.Query(tsql)

	if err != nil {
		fmt.Println("Error reading rows: " + err.Error())
		return res
	}
	for rows.Next(){
		err := rows.Scan(&item.ID, &item.PersonID, &item.UserName, &item.Password, &item.TypeUser)
		if err != nil {
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

func ValidateSystemUserLogin(user string, password string) (helper.State, string){
	items := GetSystemUserFromUserName(user)
	person := GetPatient(items[0].PersonID)
	if len(items) > 0 {
		if validatePasswordSystemUserForReset(password, items[0]){
			if len(person[0].Mail) != 0{
				newPassword := utils.CreateNewPassword()
				mail := models.Mail{
					From: person[0].Mail,
					User: user,
					Password: newPassword,
				}
				_, err := UpdatePasswordSystemUser(items[0].ID, newPassword)
				if err != nil {
					return helper.ErrorUP, ""
				}
				utils.Sendmail(mail)
				return helper.PasswordUpdate, ""
			}
			return helper.NotFoundMail, ""
		}
		if comparePassword(items[0].Password, password) {
			return helper.Accept, items[0].ID
		}
		return helper.InvalidCredentials, ""
	}
	return helper.NotFound, ""
}

func UpdatePasswordSystemUser(id string, password string) (int64, error) {
	ctx := context.Background()
	tsql := fmt.Sprintf(QuerySystemUser["updatePassword"].Q, id)
	result, err := DB.ExecContext(
		ctx,
		tsql,
		sql.Named("Password", encrypt(password)))
	if err != nil {
		return -1, err
	}
	return result.RowsAffected()
}

func validatePasswordSystemUserForReset(password string, patient models.SystemUser) bool {
	return patient.UserName == password
}

func comparePassword(hashedPassword string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return false
	}
	return true
}

func encrypt(password string) string {
	passwordByte := []byte(password)
	hashedPassword, err := bcrypt.GenerateFromPassword(passwordByte, bcrypt.DefaultCost)
	if err != nil {
		panic(err)
		return ""
	}
	return string(hashedPassword)
}


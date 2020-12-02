package db

import (
	"context"
	"crypto/md5"
	"database/sql"
	"encoding/base64"
	"fmt"
	"github.com/CarosDrean/api-results.git/constants"
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
		err := rows.Scan(&item.ID, &item.PersonID, &item.UserName, &item.Password, &item.TypeUser, &item.IsDelete)
		protocolSystemUsers := GetProtocolSystemUserWidthSystemUserID(item.ID)
		if len(protocolSystemUsers) > 0 {
			protocol := GetProtocol(protocolSystemUsers[0].ProtocolID)
			item.OrganizationID = GetOrganization(protocol.OrganizationID).ID
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
		err := rows.Scan(&item.ID, &item.PersonID, &item.UserName, &item.Password, &item.TypeUser, &item.IsDelete)
		if err != nil {
			log.Println(err)
			return res
		} else if item.IsDelete != 1{
			res = append(res, item)
		}
	}
	defer rows.Close()
	return res
}

func ValidateSystemUserLogin(user string, password string) (constants.State, string){
	items := GetSystemUserFromUserName(user)
	fmt.Println(items[0])
	if len(items) > 0 {
		if validatePasswordSystemUserForReset(password, items[0]){
			person := GetPatient(items[0].PersonID)
			if len(person[0].Mail) != 0{
				newPassword := utils.CreateNewPassword()
				mail := models.Mail{
					From: person[0].Mail,
					User: user,
					Password: newPassword,
				}
				_, err := UpdatePasswordSystemUser(items[0].ID, newPassword)
				if err != nil {
					return constants.ErrorUP, ""
				}
				utils.Sendmail(mail)
				return constants.PasswordUpdate, ""
			}
			return constants.NotFoundMail, ""
		}
		if comparePassword(items[0].Password, password) {
			return constants.Accept, items[0].ID
		}
		return constants.InvalidCredentials, ""
	}
	return constants.NotFound, ""
}

func UpdatePasswordSystemUser(id string, password string) (int64, error) {
	ctx := context.Background()
	tsql := fmt.Sprintf(QuerySystemUser["updatePassword"].Q, id)
	result, err := DB.ExecContext(
		ctx,
		tsql,
		sql.Named("Password", encryptMD5(password)))
	if err != nil {
		return -1, err
	}
	return result.RowsAffected()
}

func validatePasswordSystemUserForReset(password string, patient models.SystemUser) bool {
	return patient.UserName == password
}

func encryptMD5(text string) string {
	data := []byte(text)
	hash := md5.Sum(forSystem(data))
	return base64.StdEncoding.EncodeToString(hash[:])
}

func forSystem(data []byte) []byte {
	res := make([]byte, 0)
	for _, e := range data {
		res = append(res, e)
		res = append(res, 0)
	}
	return res
}

func comparePassword(hashedPassword string, password string) bool {
	if hashedPassword != encryptMD5(password) {
		return false
	}
	/*err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return false
	}*/
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


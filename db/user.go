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
	"strconv"
)

func GetSystemUsers() []models.SystemUser {
	res := make([]models.SystemUser, 0)
	var item models.SystemUser

	tsql := fmt.Sprintf(QuerySystemUser["list"].Q)
	rows, err := DB.Query(tsql)

	if err != nil {
		fmt.Println("Error reading rows: " + err.Error())
		return res
	}
	for rows.Next(){
		err := rows.Scan(&item.ID, &item.PersonID, &item.UserName, &item.Password, &item.TypeUser, &item.IsDelete)
		protocolSystemUsers := GetProtocolSystemUserWidthSystemUserID(strconv.FormatInt(item.ID, 10))
		if len(protocolSystemUsers) > 0 {
			protocol := GetProtocol(protocolSystemUsers[0].ProtocolID)
			item.OrganizationID = GetOrganization(protocol.OrganizationID).ID
		}

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
		protocolSystemUsers := GetProtocolSystemUserWidthSystemUserID(strconv.FormatInt(item.ID, 10))
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

func CreateSystemUser(item models.SystemUser) (int64, error) {
	ctx := context.Background()
	tsql := fmt.Sprintf(QuerySystemUser["insert"].Q)
	sequentialID := GetNextSequentialId(constants.IdNode, constants.IdSystemUserTable)
	item.Password = encryptMD5(item.Password)
	item.ID = int64(sequentialID)

	_, err := DB.ExecContext(
		ctx,
		tsql,
		sql.Named("i_SystemUserId", item.ID),
		sql.Named("v_PersonId", item.PersonID),
		sql.Named("v_UserName", item.UserName),
		sql.Named("v_Password", item.Password),
		sql.Named("i_SystemUserTypeId", item.TypeUser),
		sql.Named("i_IsDeleted", 0))
	if err != nil {
		return -1, err
	}
	return int64(sequentialID), nil
}

func UpdateSystemUser(item models.SystemUser) (int64, error) {
	ctx := context.Background()
	tsql := fmt.Sprintf(QuerySystemUser["update"].Q)
	user := GetSystemUser(strconv.FormatInt(item.ID, 10))[0]
	if user.Password != item.Password {
		item.Password = encryptMD5(item.Password)
	}

	result, err := DB.ExecContext(
		ctx,
		tsql,
		sql.Named("ID", item.ID),
		sql.Named("v_PersonId", item.PersonID),
		sql.Named("v_UserName", item.UserName),
		sql.Named("v_Password", item.Password),
		sql.Named("i_SystemUserTypeId", item.TypeUser),
		sql.Named("i_IsDeleted", 0))
	if err != nil {
		return -1, err
	}
	return result.RowsAffected()
}

func DeleteSystemUser(id string) (int64, error) {
	ctx := context.Background()
	tsql := fmt.Sprintf(QuerySystemUser["delete"].Q)
	result, err := DB.ExecContext(
		ctx,
		tsql,
		sql.Named("ID", id))
	if err != nil {
		fmt.Println(err)
		return -1, err
	}
	return result.RowsAffected()
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
			person := GetPerson(items[0].PersonID)
			if len(person[0].Mail) != 0{
				newPassword := utils.CreateNewPassword()
				mail := models.Mail{
					From: person[0].Mail,
					User: user,
					Password: newPassword,
				}
				_, err := UpdatePasswordSystemUser(strconv.FormatInt(items[0].ID, 10), newPassword)
				if err != nil {
					return constants.ErrorUP, ""
				}
				utils.SendMail(mail, constants.RouteNewPassword)
				return constants.PasswordUpdate, ""
			}
			return constants.NotFoundMail, ""
		}
		if comparePassword(items[0].Password, password) {
			return constants.Accept, strconv.FormatInt(items[0].ID, 10)
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


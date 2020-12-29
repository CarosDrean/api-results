package db

import (
	"context"
	"crypto/md5"
	"database/sql"
	"encoding/base64"
	"fmt"
	"github.com/CarosDrean/api-results.git/constants"
	"github.com/CarosDrean/api-results.git/models"
	"github.com/CarosDrean/api-results.git/query"
	"github.com/CarosDrean/api-results.git/utils"
	"golang.org/x/crypto/bcrypt"
	"strconv"
)

type UserDB struct {}

func (db UserDB) GetAll() ([]models.SystemUser, error) {
	res := make([]models.SystemUser, 0)

	tsql := fmt.Sprintf(query.SystemUser["list"].Q)
	rows, err := DB.Query(tsql)

	err = db.scan(rows, err, &res, "User DB", "Get")
	if err != nil {
		return res, err
	}
	defer rows.Close()
	return res, nil
}

func (db UserDB) Get(id string) (models.SystemUser, error) {
	res := make([]models.SystemUser, 0)

	tsql := fmt.Sprintf(query.SystemUser["get"].Q, id)
	rows, err := DB.Query(tsql)

	err = db.scan(rows, err, &res, "User DB", "Get")
	if err != nil {
		return models.SystemUser{}, err
	}
	if len(res) == 0 {
		return models.SystemUser{}, nil
	}
	defer rows.Close()
	return res[0], nil
}

func (db UserDB) Create(item models.SystemUser) (int64, error) {
	ctx := context.Background()
	tsql := fmt.Sprintf(query.SystemUser["insert"].Q)
	sequentialID := SequentialDB{}.NextSequentialId(constants.IdNode, constants.IdSystemUserTable)
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

func (db UserDB) Update(item models.SystemUser) (int64, error) {
	ctx := context.Background()
	tsql := fmt.Sprintf(query.SystemUser["update"].Q)
	user, _ := db.Get(strconv.FormatInt(item.ID, 10))
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

func (db UserDB) Delete(id string) (int64, error) {
	ctx := context.Background()
	tsql := fmt.Sprintf(query.SystemUser["delete"].Q)
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

func (db UserDB) GetFromUserName(userName string) (models.SystemUser, error) {
	res := make([]models.SystemUser, 0)
	tsql := fmt.Sprintf(query.SystemUser["getUserName"].Q, userName)
	rows, err := DB.Query(tsql)

	err = db.scan(rows, err, &res, "User DB", "Get")
	if err != nil {
		return models.SystemUser{}, err
	}
	if len(res) == 0 {
		return models.SystemUser{}, nil
	}
	defer rows.Close()
	return res[0], nil
}

func (db UserDB) ValidateLogin(user string, password string) (constants.State, string){
	item, err := db.GetFromUserName(user)
	if err != nil {
		return constants.NotFound, ""
	}
	if item.UserName == "" && item.PersonID == "" {
		return constants.NotFound, ""
	}
	if validatePasswordSystemUserForReset(password, item){
		person, _ := PersonDB{}.Get(item.PersonID)
		if len(person.Mail) != 0{
			newPassword := utils.CreateNewPassword()
			mail := models.Mail{
				From: person.Mail,
				User: user,
				Password: newPassword,
			}
			_, err := db.UpdatePassword(strconv.FormatInt(item.ID, 10), newPassword)
			if err != nil {
				return constants.ErrorUP, ""
			}
			_ = utils.SendMail(mail, constants.RouteNewPassword)
			return constants.PasswordUpdate, ""
		}
		return constants.NotFoundMail, ""
	}
	if comparePassword(item.Password, password) {
		return constants.Accept, strconv.FormatInt(item.ID, 10)
	}
	return constants.InvalidCredentials, ""

}

func (db UserDB) UpdatePassword(id string, password string) (int64, error) {
	ctx := context.Background()
	tsql := fmt.Sprintf(query.SystemUser["updatePassword"].Q, id)

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

func (db UserDB) scan(rows *sql.Rows, err error, res *[]models.SystemUser, ctx string, situation string) error {
	var item models.SystemUser
	if err != nil {
		checkError(err, situation, ctx, "Reading rows")
		return err
	}
	for rows.Next() {
		err := rows.Scan(&item.ID, &item.PersonID, &item.UserName, &item.Password, &item.TypeUser, &item.IsDelete)
		protocolSystemUsers := GetProtocolSystemUserWidthSystemUserID(strconv.FormatInt(item.ID, 10))
		if len(protocolSystemUsers) > 0 {
			protocol := GetProtocol(protocolSystemUsers[0].ProtocolID)
			organization, _ := OrganizationDB{}.Get(protocol.OrganizationID)
			item.OrganizationID = organization.ID
		}
		if err != nil {
			checkError(err, situation, ctx, "Scan rows")
			return err
		} else if item.IsDelete != 1{
			*res = append(*res, item)
		}
	}
	return nil
}


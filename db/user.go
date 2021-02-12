package db

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/CarosDrean/api-results.git/constants"
	"github.com/CarosDrean/api-results.git/models"
	"github.com/CarosDrean/api-results.git/query"
	"github.com/CarosDrean/api-results.git/utils"
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

func (db UserDB) GetAllOrganization(idOrganization string) ([]models.SystemUser, error) {
	res := make([]models.SystemUser, 0)
	var item models.SystemUser

	tsql := fmt.Sprintf(query.SystemUser["getOrganization"].Q, idOrganization)
	rows, err := DB.Query(tsql)
	if err != nil {
		checkError(err, "GetAllOrganization", "DB", "Reading rows")
		return res, err
	}

	for rows.Next() {
		err = rows.Scan(&item.TypeUser)
		if err != nil {
			checkError(err, "GetAllOrganization", "ctx", "Scan rows")
		} else if item.IsDelete != 1{
			res = append(res, item)
		}
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
	item.Password = utils.EncryptMD5(item.Password)
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
		item.Password = utils.EncryptMD5(item.Password)
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

func (db UserDB) ValidateLogin(user string, password string) (constants.State, string, error){
	item, err := db.GetFromUserName(user)
	if err != nil {
		return constants.NotFound, "", err
	}
	if item.UserName == "" && item.PersonID == "" {
		return constants.NotFound, "", nil
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
			data, _ := json.Marshal(mail)
			_, err := db.UpdatePassword(strconv.FormatInt(item.ID, 10), newPassword)
			if err != nil {
				return constants.ErrorUP, "", nil
			}
			_ = utils.SendMail(data, constants.RouteNewPassword)
			return constants.PasswordUpdate, "", nil
		}
		return constants.NotFoundMail, "", nil
	}
	if utils.ComparePassword(item.Password, password) {
		return constants.Accept, strconv.FormatInt(item.ID, 10), nil
	}
	return constants.InvalidCredentials, "", nil

}

func (db UserDB) UpdatePassword(id string, password string) (int64, error) {
	ctx := context.Background()
	tsql := fmt.Sprintf(query.SystemUser["updatePassword"].Q, id)

	result, err := DB.ExecContext(
		ctx,
		tsql,
		sql.Named("Password", utils.EncryptMD5(password)))
	if err != nil {
		return -1, err
	}
	return result.RowsAffected()
}

func validatePasswordSystemUserForReset(password string, patient models.SystemUser) bool {
	return patient.UserName == password
}

func (db UserDB) scan(rows *sql.Rows, err error, res *[]models.SystemUser, ctx string, situation string) error {
	var item models.SystemUser
	if err != nil {
		checkError(err, situation, ctx, "Reading rows")
		return err
	}
	for rows.Next() {
		err := rows.Scan(&item.ID, &item.PersonID, &item.UserName, &item.Password, &item.TypeUser, &item.IsDelete)
		protocolSystemUsers, _ := ProtocolSystemUserDB{}.GetAllSystemUserID(strconv.FormatInt(item.ID, 10))
		item.AccessClient = false
		if len(protocolSystemUsers) > 0 {
			protocol, _ := ProtocolDB{}.Get(protocolSystemUsers[0].ProtocolID)
			organization, _ := OrganizationDB{}.Get(protocol.OrganizationID)
			item.OrganizationID = organization.ID
			if protocolSystemUsers[0].ApplicationHierarchy == constants.CodeAccessClient {
				item.AccessClient = true
			}
		} else {
			item.OrganizationID = ""
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


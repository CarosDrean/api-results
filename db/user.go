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

type UserDB struct{}

func (db UserDB) GetAll() ([]models.SystemUser, error) {
	res := make([]models.SystemUser, 0)

	tsql := fmt.Sprintf(query.SystemUser["list"].Q)

	rows, err := DB.Query(tsql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		item, err := db.scanRow(rows)
		if err != nil {
			return nil, err
		}

		if item.IsDelete == 0 {
			res = append(res, item)
		}
	}

	return res, nil
}

func (db UserDB) GetAllByOrganizationID(idOrganization string) ([]models.SystemUser, error) {
	res := make([]models.SystemUser, 0)

	tsql := fmt.Sprintf(query.SystemUser["getByOrganizationID"].Q, idOrganization)

	rows, err := DB.Query(tsql)
	if err != nil {
		return res, err
	}
	defer rows.Close()

	for rows.Next() {
		item, err := db.scanRowOrganization(rows)
		if err != nil {
			return res, err
		}

		res = append(res, item)
	}

	return res, nil
}

func (db UserDB) Get(idString string) (models.SystemUser, error) {
	id, err := strconv.Atoi(idString)
	if err != nil {
		return models.SystemUser{}, err
	}

	tsql := fmt.Sprintf(query.SystemUser["get"].Q, id)

	row := DB.QueryRow(tsql)

	item, err := db.scanRow(row)
	if err != nil {
		return models.SystemUser{}, err
	}

	if item.IsDelete == 1 {
		return models.SystemUser{}, nil
	}

	return item, nil
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
		return -1, err
	}
	return result.RowsAffected()
}

func (db UserDB) GetFromUserName(userName string) (models.SystemUser, error) {
	tsql := fmt.Sprintf(query.SystemUser["getUserName"].Q, userName)

	row := DB.QueryRow(tsql)

	item, err := db.scanRow(row)
	if err != nil {
		return models.SystemUser{}, err
	}

	return item, nil
}

func (db UserDB) ValidateLogin(user string, password string, token string) (constants.State, string, error) {
	item, err := db.GetFromUserName(user)
	if err != nil {
		return constants.NotFound, "", err
	}

	if item.UserName == "" && item.PersonID == "" {
		return constants.NotFound, "", nil
	}

	if validatePasswordSystemUserForReset(password, item) {
		person, _ := PersonDB{}.Get(item.PersonID)

		if len(person.Mail) != 0 {
			newPassword := utils.CreateNewPassword()
			mail := models.Mail{
				Email:    person.Mail,
				User:     user,
				Password: newPassword,
			}

			data, _ := json.Marshal(mail)

			_, err := db.UpdatePassword(strconv.FormatInt(item.ID, 10), newPassword)
			if err != nil {
				return constants.ErrorUP, "", nil
			}

			_, err = utils.SendMail(data, constants.RouteNewPassword, token)
			if err != nil {
				return "", "", err
			}

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

func (db UserDB) scanRow(row models.RowScanner) (models.SystemUser, error) {
	var item models.SystemUser

	err := row.Scan(
		&item.ID,
		&item.PersonID,
		&item.UserName,
		&item.Password,
		&item.TypeUser,
		&item.IsDelete,
	)
	if err != nil {
		return models.SystemUser{}, err
	}

	protocolSystemUsers, _ := ProtocolSystemUserDB{}.GetAllSystemUserID(strconv.FormatInt(item.ID, 10))

	if len(protocolSystemUsers) > 0 {
		protocol, _ := ProtocolDB{}.Get(protocolSystemUsers[0].ProtocolID)
		organization, _ := OrganizationDB{}.Get(protocol.OrganizationID)

		item.OrganizationID = organization.ID

		if protocolSystemUsers[0].ApplicationHierarchy == constants.CodeAccessClient {
			item.AccessClient = true
		}
	}

	return item, nil
}

func (db UserDB) scanRowOrganization(row models.RowScanner) (models.SystemUser, error) {
	var createdAtNull sql.NullTime

	var item models.SystemUser

	err := row.Scan(
		&item.ID,
		&item.PersonID,
		&item.UserName,
		&item.Password,
		&item.TypeUser,
		&createdAtNull,
	)
	if err != nil {
		return models.SystemUser{}, err
	}
	item.CreatedAt = createdAtNull.Time

	return item, nil
}

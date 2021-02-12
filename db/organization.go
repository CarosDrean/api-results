package db

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/CarosDrean/api-results.git/models"
	"github.com/CarosDrean/api-results.git/query"
)

type OrganizationDB struct {}

func (db OrganizationDB) GetAll() ([]models.Organization, error) {
	res := make([]models.Organization, 0)

	tsql := fmt.Sprintf(query.Organization["list"].Q)
	rows, err := DB.Query(tsql)
	err = db.scan(rows, err, &res, "Organization DB", "GetAll")
	if err != nil {
		return res, err
	}
	defer rows.Close()
	return res, err
}

func (db OrganizationDB) GetAllWorkingOfEmployer(idUser string) ([]models.Organization, error) {
	res := make([]models.Organization, 0)
	var item models.Organization

	tsql := fmt.Sprintf(query.Organization["listWorkingOfEmployer"].Q, idUser)
	rows, err := DB.Query(tsql)
	if err != nil {
		checkError(err, "GetAllWorkingOfEmployer", "db", "Reading rows")
		return res, err
	}

	for rows.Next() {
		err = rows.Scan(&item.ID, &item.Name)
		if err != nil {
			checkError(err, "GetAllWorkingOfEmployer", "db", "scan rows")
		} else {
			res = append(res, item)
		}
	}
	defer rows.Close()
	return res, err
}

func (db OrganizationDB) Get(id string) (models.Organization, error) {
	res := make([]models.Organization, 0)

	tsql := fmt.Sprintf(query.Organization["get"].Q, id)
	rows, err := DB.Query(tsql)

	err = db.scan(rows, err, &res, "Organization DB", "Get")
	if err != nil {
		return models.Organization{}, err
	}
	if len(res) == 0 {
		return models.Organization{}, nil
	}
	defer rows.Close()
	return res[0], nil
}

func (db OrganizationDB) Update(id string, item models.Organization) (int64, error) {
	ctx := context.Background()
	tsql := fmt.Sprintf(query.Organization["update"].Q)
	result, err := DB.ExecContext(
		ctx,
		tsql,
		sql.Named("ID", id),
		sql.Named("v_Name", item.Name),
		sql.Named("v_Mail", item.Mail),
		sql.Named("v_EmailContacto", item.MailContact),
		sql.Named("v_EmailMedico", item.MailMedic),
		sql.Named("b_urlAdmin", item.UrlAdmin),
		sql.Named("b_urlMedic", item.UrlMedic))
	if err != nil {
		return -1, err
	}
	return result.RowsAffected()
}

func (db OrganizationDB) Delete(id string) (int64, error) {
	ctx := context.Background()
	tsql := fmt.Sprintf(query.Organization["delete"].Q)
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

func (db OrganizationDB) scan(rows *sql.Rows, err error, res *[]models.Organization, ctx string, situation string) error {
	var item models.Organization
	if err != nil {
		checkError(err, situation, ctx, "Reading rows")
		return err
	}
	for rows.Next() {
		var mailMedic sql.NullString
		var urlAdmin sql.NullBool
		var urlMedic sql.NullBool
		err := rows.Scan(&item.ID, &item.Name, &item.Mail, &item.MailContact, &mailMedic, &urlAdmin, &urlMedic)
		item.MailMedic = mailMedic.String
		item.UrlAdmin = urlAdmin.Bool
		item.UrlMedic = urlMedic.Bool
		if err != nil {
			checkError(err, situation, ctx, "Scan rows")
			return err
		} else {
			*res = append(*res, item)
		}
	}
	return nil
}

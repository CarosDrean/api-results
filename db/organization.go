package db

import (
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

func (db OrganizationDB) scan(rows *sql.Rows, err error, res *[]models.Organization, ctx string, situation string) error {
	var item models.Organization
	if err != nil {
		checkError(err, situation, ctx, "Reading rows")
		return err
	}
	for rows.Next() {
		var mailMedic sql.NullString
		err := rows.Scan(&item.ID, &item.Name, &item.Mail, &item.MailContact, &mailMedic)
		item.MailMedic = mailMedic.String
		if err != nil {
			checkError(err, situation, ctx, "Scan rows")
			return err
		} else {
			*res = append(*res, item)
		}
	}
	return nil
}

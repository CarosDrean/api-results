package db

import (
	"database/sql"
	"fmt"
	"github.com/CarosDrean/api-results.git/models"
	"github.com/CarosDrean/api-results.git/query"
	"log"
)

func GetOrganizations() []models.Organization {
	res := make([]models.Organization, 0)
	var item models.Organization

	tsql := fmt.Sprintf(query.Organization["list"].Q)
	rows, err := DB.Query(tsql)

	if err != nil {
		fmt.Println("Error reading rows: " + err.Error())
		return res
	}
	for rows.Next(){
		var mailMedic sql.NullString
		err := rows.Scan(&item.ID, &item.Name, &item.Mail, &item.MailContact, &mailMedic)
		if mailMedic.Valid {
			item.MailMedic = mailMedic.String
		} else {
			item.MailMedic = ""
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

func GetOrganization(id string) models.Organization {
	res := make([]models.Organization, 0)
	var item models.Organization

	tsql := fmt.Sprintf(query.Organization["get"].Q, id)
	rows, err := DB.Query(tsql)

	if err != nil {
		fmt.Println("Error reading rows: " + err.Error())
		return res[0]
	}
	for rows.Next(){
		var mailMedic sql.NullString
		err := rows.Scan(&item.ID, &item.Name, &item.Mail, &item.MailContact, &mailMedic)
		if mailMedic.Valid {
			item.MailMedic = mailMedic.String
		} else {
			item.MailMedic = ""
		}
		if err != nil {
			log.Println(err)
			return res[0]
		} else{
			res = append(res, item)
		}
	}
	defer rows.Close()
	return item
}

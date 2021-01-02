package db

import (
	"fmt"
	"github.com/CarosDrean/api-results.git/models"
	"github.com/CarosDrean/api-results.git/query"
	"log"
	"strings"
)

type ProtocolDB struct {}

func (db ProtocolDB) GetAllLocation(id string) []models.Protocol {
	res := make([]models.Protocol, 0)
	var item models.Protocol

	tsql := fmt.Sprintf(query.Protocol["getLocation"].Q, id)
	rows, err := DB.Query(tsql)

	if err != nil {
		fmt.Println("Error reading rows: " + err.Error())
		return res
	}
	for rows.Next(){
		err := rows.Scan(&item.ID, &item.Name, &item.OrganizationID, &item.LocationID, &item.IsDeleted, &item.EsoType)
		if err != nil {
			log.Println(err)
			return res
		} else if item.IsDeleted != 1 {
			// aqui quitar el nombre de la empresa del protocolo
			item.Name = db.delBusinessName(item.Name)
			res = append(res, item)
		}
	}
	defer rows.Close()
	return res
}

func (db ProtocolDB) GetAllOrganization(id string) []models.Protocol {
	res := make([]models.Protocol, 0)
	var item models.Protocol

	tsql := fmt.Sprintf(query.Protocol["getOrganization"].Q, id)
	rows, err := DB.Query(tsql)

	if err != nil {
		fmt.Println("Error reading rows: " + err.Error())
		return res
	}
	for rows.Next(){
		err := rows.Scan(&item.ID, &item.Name, &item.OrganizationID, &item.LocationID, &item.IsDeleted, &item.EsoType)
		if err != nil {
			log.Println(err)
			return res
		} else if item.IsDeleted != 1 {
			item.Name = db.delBusinessName(item.Name)
			res = append(res, item)
		}
	}
	defer rows.Close()
	return res
}

func (db ProtocolDB) delBusinessName(nameComplet string) string {
	pr := strings.Split(nameComplet, "-")
	name := nameComplet
	for i, e := range pr {
		if i == 1 {
			name = e
		} else if i != 0 {
			name = name + " - " + e
		}
	}
	return name
}

func (db ProtocolDB) Get(id string) models.Protocol {
	res := make([]models.Protocol, 0)
	var item models.Protocol

	tsql := fmt.Sprintf(query.Protocol["get"].Q, id)
	rows, err := DB.Query(tsql)

	if err != nil {
		fmt.Println("Error reading rows: " + err.Error())
		return res[0]
	}
	for rows.Next(){
		err := rows.Scan(&item.ID, &item.Name, &item.OrganizationID, &item.LocationID, &item.IsDeleted, &item.EsoType)
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

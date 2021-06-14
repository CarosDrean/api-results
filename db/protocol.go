package db

import (
	"database/sql"
	"fmt"
	"github.com/CarosDrean/api-results.git/models"
	"github.com/CarosDrean/api-results.git/query"
	"strings"
)

type ProtocolDB struct {}

func (db ProtocolDB) GetAllLocation(id string) ([]models.Protocol, error) {
	res := make([]models.Protocol, 0)

	tsql := fmt.Sprintf(query.Protocol["getLocation"].Q, id)
	rows, err := DB.Query(tsql)

	err = db.scan(rows, err, &res, "Protocol DB", "GetAll")
	if err != nil {
		return res, err
	}
	defer rows.Close()
	return res, nil
}
// obtener las empresas
func (db ProtocolDB) GetAllOrganization(id string) ([]models.Protocol, error) {
	res := make([]models.Protocol, 0)

	tsql := fmt.Sprintf(query.Protocol["getOrganization"].Q, id)
	rows, err := DB.Query(tsql)
	//fmt.Println(tsql)
	err = db.scan(rows, err, &res, "Protocol DB", "GetAll")
	//fmt.Println(err)
	if err != nil {
		return res, err
	}

	defer rows.Close()
	return res, nil
}
// obtener las empresas con su contratista
func (db ProtocolDB) GetAllOrganizationEmployer(id string) ([]models.Protocol, error) {
	res := make([]models.Protocol, 0)

	tsql := fmt.Sprintf(query.Protocol["getOrganizationEmployer"].Q, id)
	rows, err := DB.Query(tsql)

	err = db.scan(rows, err, &res, "Protocol DB", "GetAllOrganizationEmployer")
	if err != nil {
		return res, err
	}
	defer rows.Close()
	return res, nil
}

func (db ProtocolDB) delBusinessName(nameFull string) string {
	pr := strings.Split(nameFull, "-")
	name := nameFull
	for i, e := range pr {
		if i == 1 {
			name = e
		} else if i != 0 {
			name = name + " - " + e
		}
	}
	return name
}

func (db ProtocolDB) Get(id string) (models.Protocol, error) {
	res := make([]models.Protocol, 0)

	tsql := fmt.Sprintf(query.Protocol["get"].Q, id)
	rows, err := DB.Query(tsql)

	err = db.scan(rows, err, &res, "Location DB", "GetAll")
	if err != nil {
		return models.Protocol{}, err
	}
	if len(res) == 0 {
		return models.Protocol{}, nil
	}
	defer rows.Close()
	return res[0], nil
}


func (db ProtocolDB) scan(rows *sql.Rows, err error, res *[]models.Protocol, ctx string, situation string) error {
	var item models.Protocol
	if err != nil {
		checkError(err, situation, ctx, "Reading rows")
		return err
	}
	for rows.Next() {
		err := rows.Scan(&item.ID, &item.Name, &item.OrganizationID, &item.OrganizationEmployerID, &item.LocationID, &item.IsDeleted, &item.EsoType, &item.GroupOccupationId)
		if err != nil {
			checkError(err, situation, ctx, "Scan rows")
			return err
		} else if item.IsDeleted == 0 {
			item.BusinessName = item.Name
			item.Name = db.delBusinessName(item.Name)
			*res = append(*res, item)
		}
	}
	return nil
}

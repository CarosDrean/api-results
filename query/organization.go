package query

import "github.com/CarosDrean/api-results.git/models"

var organization = models.TableDB{
	Name:   "dbo.organization",
	Fields: []string{"v_OrganizationId", "v_Name", "v_Mail", "v_EmailContacto", "v_EmailMedico"},
}

var Organization = models.QueryDB{
	"list": {Q: "select " + fieldString(organization.Fields) + " from " + organization.Name + " order by " + organization.Fields[1] + " asc;"},
	"get":  {Q: "select " + fieldString(organization.Fields) + " from " + organization.Name + " where " + organization.Fields[0] + " = '%s';"},
}

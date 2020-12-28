package query

import "github.com/CarosDrean/api-results.git/models"

var location = models.TableDB{
	Name:   "dbo.location",
	Fields: []string{"v_LocationId", "v_OrganizationId", "v_Name", "i_IsDeleted"},
}

var Location = models.QueryDB{
	"getOrganizationID": {Q: "select " + fieldString(location.Fields) + " from " + location.Name + " where " + location.Fields[1] + " = '%s';"},
	"get":               {Q: "select " + fieldString(location.Fields) + " from " + location.Name + " where " + location.Fields[0] + " = '%s';"},
}

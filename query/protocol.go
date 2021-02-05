package query

import "github.com/CarosDrean/api-results.git/models"

var protocol = models.TableDB{
	Name:   "dbo.protocol",
	Fields: []string{"v_ProtocolId", "v_Name", "v_CustomerOrganizationId", "v_EmployerLocationId", "i_IsDeleted", "i_EsoTypeId", "v_GroupOccupationId"},
}

var Protocol = models.QueryDB{
	"getLocation":     {Q: "select " + fieldString(protocol.Fields) + " from " + protocol.Name + " where " + protocol.Fields[3] + " = '%s';"},
	"getOrganization": {Q: "select " + fieldString(protocol.Fields) + " from " + protocol.Name + " where " + protocol.Fields[2] + " = '%s';"},
	"get":             {Q: "select " + fieldString(protocol.Fields) + " from " + protocol.Name + " where " + protocol.Fields[0] + " = '%s';"},
}

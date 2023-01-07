package query

import "github.com/CarosDrean/api-results.git/models"

var protocol = models.TableDB{
	Name:   "dbo.protocol",
	Fields: []string{"v_ProtocolId", "v_Name", "v_WorkingOrganizationId", "v_EmployerOrganizationId", "v_WorkingLocationId", "i_EsoTypeId", "v_GroupOccupationId"},
}

var Protocol = models.QueryDB{
	"getLocation":             {Q: "select " + fieldString(protocol.Fields) + " from " + protocol.Name + " where " + protocol.Fields[4] + " = '%s' AND (i_IsDeleted = 0 OR i_IsDeleted IS NULL);"},
	"getOrganization":         {Q: "select " + fieldString(protocol.Fields) + " from " + protocol.Name + " where " + protocol.Fields[2] + " = '%s' AND (i_IsDeleted = 0 OR i_IsDeleted IS NULL) AND i_IsActive = 1;"},
	"getOrganizationEmployer": {Q: "select " + fieldString(protocol.Fields) + " from " + protocol.Name + " where " + protocol.Fields[3] + " = '%s' AND (i_IsDeleted = 0 OR i_IsDeleted IS NULL);"},
	"get":                     {Q: "select " + fieldString(protocol.Fields) + " from " + protocol.Name + " where " + protocol.Fields[0] + " = '%s' AND (i_IsDeleted = 0 OR i_IsDeleted IS NULL);"},
}

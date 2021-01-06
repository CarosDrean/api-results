package query

import "github.com/CarosDrean/api-results.git/models"

var organization = models.TableDB{
	Name:   "dbo.organization",
	Fields: []string{"v_OrganizationId", "v_Name", "v_Mail", "v_EmailContacto", "v_EmailMedico", "b_urlAdmin", "b_urlMedic"},
}

var Organization = models.QueryDB{
	"list":   {Q: "select " + fieldString(organization.Fields) + " from " + organization.Name + " order by " + organization.Fields[1] + " asc;"},
	"get":    {Q: "select " + fieldString(organization.Fields) + " from " + organization.Name + " where " + organization.Fields[0] + " = '%s';"},
	"update": {Q: "update " + organization.Name + " set " + updatesString(organization.Fields) + " where " + organization.Fields[0] + " = @ID;"},

	"listSystemUser": {Q: "select " + fieldStringPrefix(organization.Fields, "o") + ", su.i_SystemUserTypeId from organization o " +
		"left join protocol p on o.v_OrganizationId = p.v_CustomerOrganizationId " +
		"left join protocolsystemuser psu on p.v_ProtocolId = psu.v_ProtocolId " +
		"left join systemuser su on psu.i_SystemUserId = su.i_SystemUserId;"},
}

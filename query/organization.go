package query

import "github.com/CarosDrean/api-results.git/models"

var organization = models.TableDB{
	Name:   "dbo.organization",
	Fields: []string{"v_OrganizationId", "v_Name", "v_Mail", "v_EmailContacto", "v_EmailMedico", "b_urlAdmin", "b_urlMedic"},
}

var Organization = models.QueryDB{
	"list":   {Q: "select " + fieldString(organization.Fields) + " from " + organization.Name + " where i_IsDeleted = 0 order by " + organization.Fields[1] + " asc;"},
	"get":    {Q: "select " + fieldString(organization.Fields) + " from " + organization.Name + " where " + organization.Fields[0] + " = '%s';"},
	"update": {Q: "update " + organization.Name + " set " + updatesString(organization.Fields) + " where " + organization.Fields[0] + " = @ID;"},
	"delete": {Q: "update " + organization.Name + " set i_IsDeleted = 1 where " + user.Fields[0] + " = @ID;"},

	"listSystemUser": {Q: "select " + fieldStringPrefix(organization.Fields, "o") + ", su.i_SystemUserTypeId from organization o " +
		"left join protocol p on o.v_OrganizationId = p.v_CustomerOrganizationId " +
		"left join protocolsystemuser psu on p.v_ProtocolId = psu.v_ProtocolId " +
		"left join systemuser su on psu.i_SystemUserId = su.i_SystemUserId;"},
	"listWorkingOfEmployer": {Q: "select o.v_OrganizationId, o.v_Name from protocolsystemuser psu " +
		"inner join protocol p on psu.v_ProtocolId = p.v_ProtocolId " +
		"inner join protocol p2 on p.v_EmployerOrganizationId = p2.v_EmployerOrganizationId " +
		"inner join organization o on p2.v_WorkingOrganizationId = o.v_OrganizationId " +
		"where psu.i_SystemUserId = %s " +
		"group by o.v_OrganizationId, o.v_Name"},
 //no se usa uwu
	"getUserOrganization": {Q: "select o.v_OrganizationId, o.v_Name, pr.v_FirstName, pr.v_FirstLastName," +
		"pr.v_SecondLastName, su.i_SystemUserId, su.v_UserName, su.i_SystemUserTypeId from organization o " +
		"inner join protocol p on p.v_WorkingOrganizationId = o.v_OrganizationId " +
		"inner join protocolsystemuser psu on psu.v_ProtocolId = p.v_ProtocolId " +
		"inner join systemuser su on psu.i_SystemUserId = su.i_SystemUserId and su.i_IsDeleted != 1 " +
		"inner join person pr on su.v_PersonId = pr.v_PersonId " +
		"where o.v_OrganizationId = '%s' " +
		"group by o.v_OrganizationId, o.v_Name, pr.v_FirstName, pr.v_FirstLastName, pr.v_SecondLastName, su.i_SystemUserId, su.v_UserName, su.i_SystemUserTypeId "},
}

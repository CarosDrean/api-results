package query

import "github.com/CarosDrean/api-results.git/models"

var user = models.TableDB{
	Name:   "dbo.systemuser",
	Fields: []string{"i_SystemUserId", "v_PersonId", "v_UserName", "v_Password", "i_SystemUserTypeId", "i_IsDeleted"},
}

var SystemUser = models.QueryDB{
	"getUserName":    {Q: "select " + fieldString(user.Fields) + " from " + user.Name + " where " + user.Fields[2] + " = '%s';"},
	"get":            {Q: "select " + fieldString(user.Fields) + " from " + user.Name + " where " + user.Fields[0] + " = %s;"},
	"list":           {Q: "select " + fieldString(user.Fields) + " from " + user.Name + ";"},
	"insert":         {Q: "insert into " + user.Name + " (" + fieldString(user.Fields) + ") values (" + valuesStringNoID(user.Fields) + ");"},
	"updatePassword": {Q: "update " + user.Name + " set v_Password = @Password where " + user.Fields[0] + " = %s;"},
	"update":         {Q: "update " + user.Name + " set " + updatesString(user.Fields) + " where " + user.Fields[0] + " = @ID;"},
	"delete":         {Q: "update " + user.Name + " set " + user.Fields[5] + " = 1 where " + user.Fields[0] + " = @ID;"},

	"getOrganization": {Q: "select u.i_SystemUserTypeId from systemuser u inner join protocolsystemuser pu on u.i_SystemUserId = pu.i_SystemUserId " +
		"inner join protocol p on pu.v_ProtocolId = p.v_ProtocolId inner join person pe on u.v_PersonId = pe.v_PersonId " +
		"inner join organization o on p.v_CustomerOrganizationId = o.v_OrganizationId " +
		"where o.v_OrganizationId = '%s';"},
}

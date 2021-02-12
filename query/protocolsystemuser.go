package query

import "github.com/CarosDrean/api-results.git/models"

var protocolSystemUser = models.TableDB{
	Name:   "dbo.protocolsystemuser",
	Fields: []string{"v_ProtocolSystemUserId", "i_SystemUserId", "v_ProtocolId", "i_ApplicationHierarchyId"},
}

var ProtocolSystemUser = models.QueryDB{
	"getSystemUserID": {Q: "select " + fieldString(protocolSystemUser.Fields) + " from " + protocolSystemUser.Name + " where " + protocolSystemUser.Fields[1] + " = '%s';"},
	"get":             {Q: "select " + fieldString(protocolSystemUser.Fields) + " from " + protocolSystemUser.Name + " where " + protocolSystemUser.Fields[0] + " = '%s';"},
	"insert":          {Q: "insert into " + protocolSystemUser.Name + " (" + fieldString(protocolSystemUser.Fields) + ") values (" + valuesStringNoID(protocolSystemUser.Fields) + ");"},
	"update":          {Q: "update " + protocolSystemUser.Name + " set " + updatesString(protocolSystemUser.Fields) + " where " + protocolSystemUser.Fields[0] + " = @ID;"},
	"delete":          {Q: "delete from " + protocolSystemUser.Name + " where " + protocolSystemUser.Fields[0] + " = @ID"},
}

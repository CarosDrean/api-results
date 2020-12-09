package db

import (
	"database/sql"
)

// db es la base de datos global
var DB *sql.DB

type nameQuery string

type queryConfig struct {
	Name string
	Q    string
}

type TableDB struct {
	Name   string
	Fields []string
}

var queryResultService = map[string]*queryConfig{
	"get": {Q: "select scfv.v_Value1 from servicecomponent sc " +
		"inner join servicecomponentfields scf on sc.v_ServiceComponentId = scf.v_ServiceComponentId " +
		"inner join servicecomponentfieldvalues scfv on scf.v_ServiceComponentFieldsId = scfv.v_ServiceComponentFieldsId " +
		"where sc.v_ServiceId = '%s' and sc.v_ComponentId = '%s' and scf.v_ComponentFieldId = '%s'"},
}

var organization = TableDB{
	Name:   "dbo.organization",
	Fields: []string{"v_OrganizationId", "v_Name"},
}

var calendar = TableDB{
	Name:   "dbo.calendar",
	Fields: []string{"v_CalendarId", "v_ServiceId", "i_CalendarStatusId"},
}

var QueryCalendar = map[string]*queryConfig{
	"getServiceID": {Q: "select " + fieldString(calendar.Fields) + " from " + calendar.Name + " where " + calendar.Fields[1] + " = '%s';"},
}

var service = TableDB{
	Name: "dbo.service",
	Fields: []string{"v_ServiceId", "v_PersonId", "v_ProtocolId", "d_ServiceDate", "i_ServiceStatusId", "i_isDeleted",
		"i_AptitudeStatusId"},
}

var QueryService = map[nameQuery]*queryConfig{
	"getPersonID": {Q: "select " + fieldString(service.Fields) + " from " + service.Name + " where " + service.Fields[1] + " = '%s' order by " + service.Fields[3] + " desc;"},
	"getProtocol": {Q: "select " + fieldString(service.Fields) + " from " + service.Name + " where " + service.Fields[2] +
		" = '%s' and d_ServiceDate is not null order by " + service.Fields[3] + " desc;"},
	"getProtocolFilter": {Q: "select " + fieldString(service.Fields) + " from " + service.Name + " where " + service.Fields[2] +
		" = '%s' and " + service.Fields[3] + ">= CONVERT(DATETIME, '%s', 102) and " + service.Fields[3] + "<= CONVERT(DATETIME, '%s', 102) and " + service.Fields[3] +
		" is not null order by " + service.Fields[3] + " desc;"},
	"get": {Q: "select " + fieldString(service.Fields) + " from " + service.Name + " where " + service.Fields[0] + " = '%s';"},
}

var protocol = TableDB{
	Name:   "dbo.protocol",
	Fields: []string{"v_ProtocolId", "v_Name", "v_CustomerOrganizationId", "v_EmployerLocationId", "i_IsDeleted"},
}

var QueryProtocol = map[string]*queryConfig{
	"getLocation":     {Q: "select " + fieldString(protocol.Fields) + " from " + protocol.Name + " where " + protocol.Fields[3] + " = '%s';"},
	"getOrganization": {Q: "select " + fieldString(protocol.Fields) + " from " + protocol.Name + " where " + protocol.Fields[2] + " = '%s';"},
	"get":             {Q: "select " + fieldString(protocol.Fields) + " from " + protocol.Name + " where " + protocol.Fields[0] + " = '%s';"},
}

var QueryOrganization = map[string]*queryConfig{
	"get": {Q: "select " + fieldString(organization.Fields) + " from " + organization.Name + " where " + organization.Fields[0] + " = '%s';"},
}

var patient = TableDB{
	Name: "dbo.person",
	Fields: []string{"v_PersonId", "v_DocNumber", "v_Password", "v_FirstName", "v_FirstLastName", "v_SecondLastName",
		"v_Mail", "i_SexTypeId", "d_Birthdate"},
}

var protocolSystemUser = TableDB{
	Name:   "dbo.protocolsystemuser",
	Fields: []string{"v_ProtocolSystemUserId", "i_SystemUserId", "v_ProtocolId"},
}

var QueryProtocolSystemUser = map[string]*queryConfig{
	"getSystemUserID": {Q: "select " + fieldString(protocolSystemUser.Fields) + " from " + protocolSystemUser.Name + " where " + protocolSystemUser.Fields[1] + " = '%s';"},
	"get":             {Q: "select " + fieldString(protocolSystemUser.Fields) + " from " + protocolSystemUser.Name + " where " + protocolSystemUser.Fields[0] + " = '%s';"},
}

var location = TableDB{
	Name:   "dbo.location",
	Fields: []string{"v_LocationId", "v_OrganizationId", "v_Name", "i_IsDeleted"},
}

var queryLocation = map[string]*queryConfig{
	"getOrganizationID": {Q: "select " + fieldString(location.Fields) + " from " + location.Name + " where " + location.Fields[1] + " = '%s';"},
	"get":               {Q: "select " + fieldString(location.Fields) + " from " + location.Name + " where " + location.Fields[0] + " = '%s';"},
}

var user = TableDB{
	Name:   "dbo.systemuser",
	Fields: []string{"i_SystemUserId", "v_PersonId", "v_UserName", "v_Password", "i_SystemUserTypeId", "i_IsDeleted"},
}

var QuerySystemUser = map[string]*queryConfig{
	"getUserName":    {Q: "select " + fieldString(user.Fields) + " from " + user.Name + " where " + user.Fields[2] + " = '%s';"},
	"get":            {Q: "select " + fieldString(user.Fields) + " from " + user.Name + " where " + user.Fields[0] + " = %s;"},
	"list":           {Q: "select " + fieldString(user.Fields) + " from " + user.Name + ";"},
	"updatePassword": {Q: "update " + user.Name + " set v_Password = @Password where " + user.Fields[0] + " = %s;"},
}

var QueryPatient = map[string]*queryConfig{
	"get":            {Q: "select " + fieldString(patient.Fields) + " from " + patient.Name + " where " + patient.Fields[0] + " = '%s';"},
	"getDNI":         {Q: "select " + fieldString(patient.Fields) + " from " + patient.Name + " where " + patient.Fields[1] + " = '%s';"},
	"list":           {Q: "select " + fieldString(patient.Fields) + " from " + patient.Name + ";"},
	"insert":         {Q: "insert into (" + fieldString(patient.Fields) + ") values (" + valuesString(patient.Fields) + ");"},
	"update":         {Q: "update " + patient.Name + " set " + updatesString(patient.Fields) + " where " + patient.Fields[0] + " = '%s';"},
	"updatePassword": {Q: "update " + patient.Name + " set v_Password = @Password where " + patient.Fields[0] + " = '%s';"},
	"delete":         {Q: "delete from " + patient.Name + " where " + patient.Fields[0] + " = ?;"},
}

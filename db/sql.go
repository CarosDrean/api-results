package db

import (
	"database/sql"
)

// db es la base de datos global
var DB *sql.DB

type queryConfig struct {
	Name string
	Q    string
}

type TableDB struct {
	Name   string
	Fields []string
}

func fieldString(fields []string) string {
	fieldString := ""
	for i, field := range fields {
		if i == 0 {
			fieldString = field
		} else {
			fieldString = fieldString + ", " + field
		}
	}
	return fieldString
}

func valuesString(fields []string) string {
	values := ""
	for i := range fields {
		if i == 0 {
			values = "?"
		} else {
			values = values + ", ?"
		}
	}
	return values
}

func updatesString(fields []string) string {
	values := ""
	for i, field := range fields {
		if i == 1 {
			values = field + " = @" + field
		} else if i != 0 {
			values = values + ", " + field + " = @" + field
		}
	}
	return values
}

var service = TableDB{
	Name:   "dbo.service",
	Fields: []string{"v_ServiceId", "v_PersonId", "v_ProtocolId", "d_ServiceDate", "i_ServiceStatusId", "i_isDeleted"},
}

var protocol = TableDB{
	Name:   "dbo.protocol",
	Fields: []string{"v_ProtocolId", "v_Name", "v_CustomerOrganizationId"},
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

var QueryService = map[string]*queryConfig{
	"getPersonID": {Q: "select " + fieldString(service.Fields) + " from " + service.Name + " where " + service.Fields[1] + " = '%s' order by " + service.Fields[3] + " desc;"},
	"get": {Q: "select " + fieldString(service.Fields) + " from " + service.Name + " where " + service.Fields[0] + " = '%s';"},
}

var QueryProtocol = map[string]*queryConfig{
	"get": {Q: "select " + fieldString(protocol.Fields) + " from " + protocol.Name + " where " + protocol.Fields[0] + " = '%s';"},
}

var QueryOrganization = map[string]*queryConfig{
	"get": {Q: "select " + fieldString(organization.Fields) + " from " + organization.Name + " where " + organization.Fields[0] + " = '%s';"},
}

var patient = TableDB{
	Name:   "dbo.person",
	Fields: []string{"v_PersonId", "v_DocNumber", "v_Password", "v_FirstName", "v_FirstLastName", "v_SecondLastName", "v_Mail"},
}

var protocolSystemUser = TableDB{
	Name:   "dbo.protocolsystemuser",
	Fields: []string{"v_ProtocolSystemUserId", "i_SystemUserId", "v_ProtocolId"},
}

var QueryProtocolSystemUser = map[string]*queryConfig{
	"getSystemUserID": {Q: "select " + fieldString(protocolSystemUser.Fields) + " from " + protocolSystemUser.Name + " where " + protocolSystemUser.Fields[1] + " = '%s';"},
	"get": {Q: "select " + fieldString(protocolSystemUser.Fields) + " from " + protocolSystemUser.Name + " where " + protocolSystemUser.Fields[0] + " = '%s';"},
}

var location = TableDB{
	Name:   "dbo.location",
	Fields: []string{"v_LocationId", "v_OrganizationId", "v_Name", "i_IsDelete"},
}

var queryLocation = map[string]*queryConfig{
	"getOrganizationID": {Q: "select " + fieldString(location.Fields) + " from " + location.Name + " where " + location.Fields[1] + " = '%s';"},
	"get": {Q: "select " + fieldString(location.Fields) + " from " + location.Name + " where " + location.Fields[0] + " = '%s';"},
}

var user = TableDB{
	Name:   "dbo.systemuser",
	Fields: []string{"i_SystemUserId", "v_PersonId", "v_UserName", "v_Password", "i_SystemUserTypeId", "i_IsDeleted"},
}

var QuerySystemUser = map[string]*queryConfig{
	"getUserName": {Q: "select " + fieldString(user.Fields) + " from " + user.Name + " where " + user.Fields[2] + " = '%s';"},
	"get": {Q: "select " + fieldString(user.Fields) + " from " + user.Name + " where " + user.Fields[0] + " = %s;"},
	"updatePassword": {Q: "update " + user.Name + " set v_Password = @Password where " + user.Fields[0] + " = %s;"},
}

var QueryPatient = map[string]*queryConfig{
	"get":    {Q: "select " + fieldString(patient.Fields) + " from " + patient.Name + " where " + patient.Fields[0] + " = '%s';"},
	"getDNI": {Q: "select " + fieldString(patient.Fields) + " from " + patient.Name + " where " + patient.Fields[1] + " = '%s';"},
	"list":   {Q: "select " + fieldString(patient.Fields) + " from " + patient.Name + ";"},
	"insert": {Q: "insert into (" + fieldString(patient.Fields) + ") values (" + valuesString(patient.Fields) + ");"},
	"update": {Q: "update " + patient.Name + " set " + updatesString(patient.Fields) + " where " + patient.Fields[0] + " = '%s';"},
	"updatePassword": {Q: "update " + patient.Name + " set v_Password = @Password where " + patient.Fields[0] + " = '%s';"},
	"delete": {Q: "delete from " + patient.Name + " where " + patient.Fields[0] + " = ?;"},
}

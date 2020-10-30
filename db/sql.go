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

var user = TableDB{
	Name:   "dbo.systemuser",
	Fields: []string{"i_SystemUserId", "v_UserName", "v_Password"},
}

var service = TableDB{
	Name:   "dbo.service",
	Fields: []string{"v_ServiceId", "v_PersonId", "v_ProtocolId", "d_ServiceDate"},
}

var protocol = TableDB{
	Name:   "dbo.protocol",
	Fields: []string{"v_ProtocolId", "v_CustomerOrganizationId"},
}

var organization = TableDB{
	Name:   "dbo.organization",
	Fields: []string{"v_OrganizationId", "v_Name"},
}

var QueryService = map[string]*queryConfig{
	"getPersonID": {Q: "select " + fieldString(service.Fields) + " from " + service.Name + " where " + service.Fields[1] + " = '%s';"},
}

var QueryProtocol = map[string]*queryConfig{
	"get": {Q: "select " + fieldString(protocol.Fields) + " from " + protocol.Name + " where " + protocol.Fields[0] + " = '%s';"},
}

var QueryOrganization = map[string]*queryConfig{
	"get": {Q: "select " + fieldString(organization.Fields) + " from " + organization.Name + " where " + organization.Fields[0] + " = '%s';"},
}

var patient = TableDB{
	Name:   "dbo.person",
	Fields: []string{"v_PersonId", "v_DocNumber", "v_Password", "v_FirstName", "v_FirstLastName", "v_SecondLastName"},
}

var QueryPatient = map[string]*queryConfig{
	"get":    {Q: "select " + fieldString(patient.Fields) + " from " + patient.Name + " where " + patient.Fields[0] + " = ?;"},
	"getDNI": {Q: "select " + fieldString(patient.Fields) + " from " + patient.Name + " where " + patient.Fields[1] + " = '%s';"},
	"list":   {Q: "select " + fieldString(patient.Fields) + " from " + patient.Name + ";"},
	"insert": {Q: "insert into (" + fieldString(patient.Fields) + ") values (" + valuesString(patient.Fields) + ");"},
	"update": {Q: "update " + patient.Name + " set " + updatesString(patient.Fields) + " where " + patient.Fields[0] + " = '%s';"},
	"updatePassword": {Q: "update " + patient.Name + " set v_Password = @Password where " + patient.Fields[0] + " = '%s';"},
	"delete": {Q: "delete from " + patient.Name + " where " + patient.Fields[0] + " = ?;"},
}

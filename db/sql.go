package db

import (
	"database/sql"
)

// db es la base de datos global
var DB *sql.DB

type stmtConfig struct {
	Stmt *sql.Stmt
	Q    string
}

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
		if i == 0 {
			values = field + " = ?"
		} else {
			values = values + ", " + field + " = ?"
		}
	}
	return values
}

var user = TableDB{
	Name:   "dbo.systemuser",
	Fields: []string{"i_SystemUserId", "v_UserName", "v_Password"},
}

var PrepStmtsUser = map[string]*stmtConfig{
	"get":    {Q: "select " + fieldString(user.Fields) + " from " + user.Name + " where " + user.Fields[0] + " = ?;"},
	"list":   {Q: "select " + fieldString(user.Fields) + " from " + user.Name + ";"},
	"insert": {Q: "insert into (" + fieldString(user.Fields) + ") values (" + valuesString(user.Fields) + ");"},
	"update": {Q: "update " + user.Name + " set " + updatesString(user.Fields) + " where " + user.Fields[0] + " = ?;"},
	"delete": {Q: "delete from " + user.Name + " where " + user.Fields[0] + " = ?;"},
}

var patient = TableDB{
	Name:   "dbo.person",
	Fields: []string{"v_PersonId", "v_DocNumber", "v_Password"},
}

var QueryPatient = map[string]*queryConfig{
	"get":    {Q: "select " + fieldString(patient.Fields) + " from " + patient.Name + " where " + patient.Fields[0] + " = ?;"},
	"getDNI": {Q: "select " + fieldString(patient.Fields) + " from " + patient.Name + " where " + patient.Fields[1] + " = '%s';"},
	"list":   {Q: "select " + fieldString(patient.Fields) + " from " + patient.Name + ";"},
	"insert": {Q: "insert into (" + fieldString(patient.Fields) + ") values (" + valuesString(patient.Fields) + ");"},
	"update": {Q: "update " + patient.Name + " set " + updatesString(patient.Fields) + " where " + patient.Fields[0] + " = ?;"},
	"delete": {Q: "delete from " + patient.Name + " where " + patient.Fields[0] + " = ?;"},
}

var PrepStmtsPatient = map[string]*stmtConfig{
	"get":    {Q: "select " + fieldString(patient.Fields) + " from " + patient.Name + " where " + patient.Fields[0] + " = ?;"},
	"getDNI": {Q: "select " + fieldString(patient.Fields) + " from " + patient.Name + " where " + patient.Fields[1] + " = '?';"},
	"list":   {Q: "select " + fieldString(patient.Fields) + " from " + patient.Name + ";"},
	"insert": {Q: "insert into (" + fieldString(patient.Fields) + ") values (" + valuesString(patient.Fields) + ");"},
	"update": {Q: "update " + patient.Name + " set " + updatesString(patient.Fields) + " where " + patient.Fields[0] + " = ?;"},
	"delete": {Q: "delete from " + patient.Name + " where " + patient.Fields[0] + " = ?;"},
}

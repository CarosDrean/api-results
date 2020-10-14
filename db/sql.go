package db

import (
	"database/sql"
)

// db es la base de datos global
var DB *sql.DB

// Prepared statements
type stmtConfig struct {
	Stmt *sql.Stmt
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
	"insert": {Q: "insert into post (" + fieldString(user.Fields) + ") values (" + valuesString(user.Fields) + ");"},
	"update": {Q: "update " + user.Name + " set " + updatesString(user.Fields) + " where " + user.Fields[0] + " = ?|;"},
	"delete": {Q: "delete from " + user.Name + " where " + user.Fields[0] + " = ?;"},
}
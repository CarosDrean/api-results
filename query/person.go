package query

import "github.com/CarosDrean/api-results.git/models"

var person = models.TableDB{
	Name: "dbo.person",
	Fields: []string{"v_PersonId", "v_DocNumber", "v_Password", "v_FirstName", "v_FirstLastName", "v_SecondLastName",
		"v_Mail", "i_SexTypeId", "d_Birthdate", "i_IsDeleted"},
}

var Person = models.QueryDB{
	"get":            {Q: "select " + fieldString(person.Fields) + " from " + person.Name + " where " + person.Fields[0] + " = '%s';"},
	"getDNI":         {Q: "select " + fieldString(person.Fields) + " from " + person.Name + " where " + person.Fields[1] + " = '%s';"},
	"list":           {Q: "select " + fieldString(person.Fields) + " from " + person.Name + ";"},
	"insert":         {Q: "insert into " + person.Name + " (" + fieldString(person.Fields) + ") values (" + valuesStringNoID(person.Fields) + ");"},
	"update":         {Q: "update " + person.Name + " set " + updatesString(person.Fields) + " where " + person.Fields[0] + " = '%s';"},
	"updatePassword": {Q: "update " + person.Name + " set v_Password = @Password where " + person.Fields[0] + " = '%s';"},
}

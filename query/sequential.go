package query

import "github.com/CarosDrean/api-results.git/models"

var sequential = models.TableDB{
	Name:   "dbo.secuential",
	Fields: []string{"i_NodeId", "i_TableId", "i_SecuentialId"},
}

var Sequential = models.QueryDB{
	"get":    {Q: "select " + fieldString(sequential.Fields) + " from " + sequential.Name + " where " + sequential.Fields[0] + " = %d and " + sequential.Fields[1] + " = %d;"},
	"insert": {Q: "insert into " + sequential.Name + " (" + fieldString(sequential.Fields) + ") values (" + valuesStringNoID(sequential.Fields) + ");"},
	"update": {Q: "update " + sequential.Name + " set " + updatesStringNoID(sequential.Fields) + " where " + sequential.Fields[0] + " = %d and " + sequential.Fields[1] + " = %d;"},
}

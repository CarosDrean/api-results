package query

import "github.com/CarosDrean/api-results.git/models"

var cie10 = models.TableDB{
	Name:   "dbo.cie10",
	Fields: []string{"v_CIE10Id", "v_CIE10Description1", "v_CIE10Description1"},
}

var CIE10 = models.QueryDB{
	"list":   {Q: "select " + fieldString(cie10.Fields) + " from " + cie10.Name + ";"},
	"insert": {Q: "insert into " + cie10.Name + " (" + fieldString(cie10.Fields) + ") values (" + valuesStringNoID(cie10.Fields) + ");"},
	"update": {Q: "update " + cie10.Name + " set " + updatesString(cie10.Fields) + " where " + cie10.Fields[0] + " = @ID;"},
}

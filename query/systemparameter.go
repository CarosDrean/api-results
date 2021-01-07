package query

import "github.com/CarosDrean/api-results.git/models"

var spT = models.TableDB{
	Name:   "dbo.systemparameter",
	Fields: []string{"i_GroupId", "i_ParameterId", "v_Value1", "i_IsDeleted"},
}

var SystemParameter = models.QueryDB{
	"getGroup": {Q: "select " + fieldString(spT.Fields) + " from " + spT.Name + " where " + spT.Fields[0] + " = %s;"},
	"list":     {Q: "select " + fieldString(spT.Fields) + " from " + spT.Name + ";"},
	"insert":   {Q: "insert into " + spT.Name + " (" + fieldString(spT.Fields) + ") values (" + valuesStringNoID(spT.Fields) + ");"},
	"update":   {Q: "update " + spT.Name + " set " + updatesString(spT.Fields) + " where " + spT.Fields[0] + " = @ID and " + spT.Fields[1] + " = IDT;"},
	"delete":   {Q: "update " + spT.Name + " set " + spT.Fields[3] + " = 1 where " + spT.Fields[0] + " = @ID and " + spT.Fields[1] + " = IDT;"},
}

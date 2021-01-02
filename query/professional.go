package query

import "github.com/CarosDrean/api-results.git/models"

var professional = models.TableDB{
	Name:   "dbo.professional",
	Fields: []string{"v_PersonId", "i_ProfessionId", "v_ProfessionalCode", "i_IsDeleted"},
}

var Professional = models.QueryDB{
	"list":   {Q: "select " + fieldString(professional.Fields) + " from " + professional.Name + ";"},
	"get":    {Q: "select " + fieldString(professional.Fields) + " from " + professional.Name + " where " + professional.Fields[0] + " = '%s';"},
	"insert": {Q: "insert into " + professional.Name + " (" + fieldString(professional.Fields) + ") values (" + valuesStringNoID(professional.Fields) + ");"},
	"update": {Q: "update " + professional.Name + " set " + updatesString(professional.Fields) + " where " + professional.Fields[0] + " = @ID;"},
	"delete": {Q: "update " + user.Name + " set " + user.Fields[5] + " = 1 where " + user.Fields[3] + " = @ID;"},
}

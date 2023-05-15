package query

import "github.com/CarosDrean/api-results.git/models"

var disease = models.TableDB{
	Name:   "dbo.diseases",
	Fields: []string{"v_DiseasesId", "v_CIE10Id", "v_Name", "i_IsDeleted"},
}

var Disease = models.QueryDB{
	"list":   {Q: "select " + fieldString(disease.Fields) + " from " + disease.Name + ";"},
	"insert": {Q: "insert into " + disease.Name + " (" + fieldString(disease.Fields) + ") values (" + valuesStringNoID(disease.Fields) + ");"},
	"update": {Q: "update " + disease.Name + " set " + updatesString(disease.Fields) + " where " + disease.Fields[0] + " = @ID;"},
	"delete": {Q: "update " + disease.Name + " set " + disease.Fields[3] + " = 1 where " + disease.Fields[0] + " = @ID;"},
}

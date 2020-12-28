package query

import "github.com/CarosDrean/api-results.git/models"

var systemParameter = models.TableDB{
	Name:   "dbo.systemparameter",
	Fields: []string{"i_GroupId", "i_ParameterId", "v_Value1"},
}

var SystemParameter = models.QueryDB{
	"getGroup": {Q: "select " + fieldString(systemParameter.Fields) + " from " + systemParameter.Name + " where " + calendar.Fields[0] + " = %s;"},
}

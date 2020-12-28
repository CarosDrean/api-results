package query

import "github.com/CarosDrean/api-results.git/models"

var component = models.TableDB{
	Name:   "dbo.component",
	Fields: []string{"v_ComponentId", "v_Name", "i_CategoryId", "i_IsDeleted"},
}

var Component = models.QueryDB{
	"getCategory": {Q: "select " + fieldString(component.Fields) + " from " + component.Name + " where " + calendar.Fields[2] + " = %s;"},
}

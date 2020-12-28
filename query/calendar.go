package query

import "github.com/CarosDrean/api-results.git/models"

var calendar = models.TableDB{
	Name:   "dbo.calendar",
	Fields: []string{"v_CalendarId", "v_ServiceId", "i_CalendarStatusId"},
}

var Calendar = models.QueryDB{
	"getServiceID": {Q: "select " + fieldString(calendar.Fields) + " from " + calendar.Name + " where " + calendar.Fields[1] + " = '%s';"},
}
package query

import "github.com/CarosDrean/api-results.git/models"

var calendar = models.TableDB{
	Name:   "dbo.calendar",
	Fields: []string{"v_CalendarId", "v_PersonId", "v_ServiceId", "v_ProtocolId", "i_CalendarStatusId", "d_CircuitStartDate", "d_SalidaCM", "d_DateTimeCalendar"},
}

var Calendar = models.QueryDB{
	"getServiceID":     {Q: "select " + fieldString(calendar.Fields) + " from " + calendar.Name + " where " + calendar.Fields[2] + " = '%s';"},
	"listOrganization": {Q: ""},
	"get":              {Q: ""},
	"create":           {Q: ""},
	"updated":          {Q: ""},
	"delete":           {Q: ""},
}

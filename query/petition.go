package query

import "github.com/CarosDrean/api-results.git/models"

var petit = models.TableDB{
	Name:   "dbo.calendarPetition",
	Fields: []string{"v_PersonId", "i_DocTypeId", "v_DocNumber", "v_FirstName", "v_FirstLastName", "v_SecondLastName",
			"i_SexTypeId", "d_Birthdate", "v_TelephoneNumber", "v_CurrentOccupation", "d_DateProgramming", "i_ServiceTypeId",
			"v_PersonProgramming", "v_ResponsableProgramming", "v_CalendarId_2", "v_WorkersCondition", "v_FactCR",
			"v_NombreProyecto", "v_OrganizationId", "v_ProtocolId", "d_deleted", "v_PetitionStatus", "v_Comentary"},
}

var Petition = models.QueryDB{
	"insert":         {Q: "insert into " + petit.Name + " (" + fieldString(petit.Fields) + ", d_DateInsert) values (" + valuesStringNoID(petit.Fields) + ", GETDATE());"},
}
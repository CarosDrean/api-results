package query

import "github.com/CarosDrean/api-results.git/models"

var petit = models.TableDB{
	Name: "dbo.calendarPetition",
	Fields: []string{"v_PersonId", "i_DocTypeId", "v_DocNumber", "v_FirstName", "v_FirstLastName", "v_SecondLastName",
		"i_SexTypeId", "d_Birthdate", "v_TelephoneNumber", "v_CurrentOccupation", "d_DateProgramming", "i_ServiceTypeId",
		"v_PersonProgramming", "v_ResponsableProgramming", "v_CalendarId_2", "v_WorkersCondition", "v_FactCR",
		"v_NombreProyecto", "v_OrganizationId", "v_ProtocolId", "d_deleted", "v_PetitionStatus", "v_Comentary"},
}

var cita = models.TableDB{
	Name: "dbo.ProcedimientosCardiologicos",
	Fields: []string{"v_Name", "v_ApePaterno", "v_ApeMaterno", "v_Doc", "v_email", "v_telefono", "v_direccion", "v_dob",
		"v_fechaConsulta", "v_procedimiento", "v_mensaje", "v_sex"},
}

var Petition = models.QueryDB{
	"insert": {Q: "insert into " + petit.Name + " (" + fieldString(petit.Fields) + ", d_DateInsert) values (" + valuesStringNoID(petit.Fields) + ", GETDATE());"},
}

var Citas = models.QueryDB{
	"insert": {Q: "insert into " + cita.Name + " (" + fieldString(cita.Fields) + ", v_fechaInsercion) values (" + valuesStringNoID(cita.Fields) + ", GETDATE());"},
}

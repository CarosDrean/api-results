package db

import (
	"context"
	_ "context"
	"database/sql"
	_ "database/sql"
	_ "encoding/json"
	"fmt"
	_ "fmt"
	_ "github.com/CarosDrean/api-results.git/constants"
	"github.com/CarosDrean/api-results.git/models"
	_ "github.com/CarosDrean/api-results.git/models"
	"github.com/CarosDrean/api-results.git/query"
	_ "github.com/CarosDrean/api-results.git/query"
	_ "github.com/CarosDrean/api-results.git/utils"
	_ "strconv"
)

type PetitionDB struct{}

func (db PetitionDB) Create(item models.PetitionProgrammation) (int64, error) {
	ctx := context.Background()

	tsql := fmt.Sprintf(query.Petition["insert"].Q)

	result, err := DB.ExecContext(
		ctx,
		tsql,
		sql.Named("v_PersonId", item.PersonId),
		sql.Named("i_DocTypeId", item.DocType),
		sql.Named("v_DocNumber", item.DocNumber),
		sql.Named("v_FirstName", item.FirstName),
		sql.Named("v_FirstLastName", item.FirstLastName),
		sql.Named("v_SecondLastName", item.SecondLastName),
		sql.Named("i_SexTypeId", item.SexTypeId),
		sql.Named("d_Birthdate", item.Birthdate),
		sql.Named("v_TelephoneNumber", item.TelephoneNumber),
		sql.Named("v_CurrentOccupation", item.CurrentOccupation),
		sql.Named("d_DateProgramming", item.DateProgramming),
		sql.Named("i_ServiceTypeId", item.ServiceTypeId),
		sql.Named("v_PersonProgramming", item.PersonProgramming),
		sql.Named("v_ResponsableProgramming", item.ResponsableProgramming),
		sql.Named("v_CalendarId_2", item.CalendarId_2),
		sql.Named("v_WorkersCondition", item.WorkersCondition),
		sql.Named("v_FactCR", item.FactCR),
		sql.Named("v_NombreProyecto", item.NombreProyecto),
		sql.Named("v_OrganizationId", item.OrganizationId),
		sql.Named("v_ProtocolId", item.ProtocolId),
		sql.Named("d_deleted", item.Deleted),
		sql.Named("v_PetitionStatus", item.PetitionStatus),
		sql.Named("v_Comentary", item.Comentary))
	if err != nil {
		return -1, err
	}
	return result.RowsAffected()
}
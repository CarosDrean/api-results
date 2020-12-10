package controller

import (
	"encoding/json"
	"github.com/CarosDrean/api-results.git/constants"
	"github.com/CarosDrean/api-results.git/db"
	"github.com/CarosDrean/api-results.git/models"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"strings"
)

func GetExams(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	var params = mux.Vars(r)
	id, _ := params["id"]
	res := make([]models.Result, 0)
	var item models.Result

	patients := db.GetPatient(id)
	services := db.GetService(patients[0].ID, db.NQGetServicePersonID)

	for i, e := range services {
		if e.ServiceStatusId == 3 && e.IsDeleted != 1 { // culminado
			calendar := db.GetCalendarService(e.ID)
			if calendar.CalendarStatusID != 4 { // 4 = cancelado
				item.ID = strconv.Itoa(i)
				item.ServiceDate = e.ServiceDate
				item.IdService = e.ID
				item.ProtocolName = db.GetProtocol(e.ProtocolID).Name
				be := strings.FieldsFunc(item.ProtocolName, Split)
				item.Business = be[0]
				item.Exam = be[len(be)-1]
				item.Result = db.GetResultService(e.ID, constants.IdPruebaRapida, constants.IdResultPruebaRapida)
				res = append(res, item)
			}
		}

	}
	_ = json.NewEncoder(w).Encode(res)
}

func Split(r rune) bool {
	return r == '-' || r == '/'
}

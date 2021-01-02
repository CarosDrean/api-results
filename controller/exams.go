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

type ExamController struct {}

func (c ExamController) GetAllPerson(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	var params = mux.Vars(r)
	id, _ := params["id"]
	res := make([]models.Result, 0)
	var item models.Result

	person, _ := db.PersonDB{}.Get(id)
	services, _ := db.ServiceDB{}.GetAllPerson(person.ID)

	for i, e := range services {
		if e.ServiceStatusId == 3 && e.IsDeleted != 1 { // culminado
			calendar, _ := db.CalendarDB{}.GetService(e.ID)
			if calendar.CalendarStatusID != 4 { // 4 = cancelado
				item.ID = strconv.Itoa(i)
				item.ServiceDate = e.ServiceDate
				item.IdService = e.ID
				item.ProtocolName = db.ProtocolDB{}.Get(e.ProtocolID).Name
				be := strings.FieldsFunc(item.ProtocolName, c.split)
				item.Business = be[0]
				item.Exam = be[len(be)-1]
				item.Result = db.GetResultService(e.ID, constants.IdPruebaRapida, constants.IdResultPruebaRapida)
				res = append(res, item)
			}
		}

	}
	_ = json.NewEncoder(w).Encode(res)
}

func (c ExamController) split(r rune) bool {
	return r == '-' || r == '/'
}

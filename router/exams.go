package routes

import (
	exam "github.com/CarosDrean/api-results.git/controller"
	mid "github.com/CarosDrean/api-results.git/middleware"
	"github.com/gorilla/mux"
)

func examsRoutes(s *mux.Router) {
	ctrl := exam.ExamController{}
	s.HandleFunc("/all/{id}", mid.CheckSecurity(ctrl.GetAllPerson)).Methods("GET")
}

package routes

import (
	exam "github.com/CarosDrean/api-results.git/controller"
	mid "github.com/CarosDrean/api-results.git/middleware"
	"github.com/gorilla/mux"
)

func examsRoutes(s *mux.Router) {
	s.HandleFunc("/all/{id}", mid.CheckSecurity(exam.GetExams)).Methods("GET")
}

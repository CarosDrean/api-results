package routes

import (
	systemParameter "github.com/CarosDrean/api-results.git/controller"
	mid "github.com/CarosDrean/api-results.git/middleware"
	"github.com/gorilla/mux"
)

func systemParameterRoutes(s *mux.Router) {
	ctrl := systemParameter.SystemParameterController{}
	s.HandleFunc("/consultings", mid.CheckSecurity(ctrl.GetConsultingS)).Methods("GET")
	s.HandleFunc("/", mid.CheckSecurity(ctrl.Create)).Methods("POST")
	s.HandleFunc("/{id}", mid.CheckSecurity(ctrl.Update)).Methods("PUT")
	s.HandleFunc("/{id}", mid.CheckSecurity(ctrl.Delete)).Methods("DELETE")
}

package routes

import (
	systemParameter "github.com/CarosDrean/api-results.git/controller"
	mid "github.com/CarosDrean/api-results.git/middleware"
	"github.com/gorilla/mux"
)

func systemParameterRoutes(s *mux.Router) {
	s.HandleFunc("/consultings", mid.CheckSecurity(systemParameter.GetConsultings)).Methods("GET")
}
